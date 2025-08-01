package template

import (
	"context"
	"fmt"
	"sync"

	"github.com/e2b-dev/infra/packages/orchestrator/internal/sandbox/block"
	blockmetrics "github.com/e2b-dev/infra/packages/orchestrator/internal/sandbox/block/metrics"
	"github.com/e2b-dev/infra/packages/orchestrator/internal/sandbox/build"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/header"
	"github.com/e2b-dev/infra/packages/shared/pkg/utils"
)

type storageTemplate struct {
	files storage.TemplateCacheFiles

	memfile  *utils.SetOnce[block.ReadonlyDevice]
	rootfs   *utils.SetOnce[block.ReadonlyDevice]
	snapfile *utils.SetOnce[File]

	memfileHeader *header.Header
	rootfsHeader  *header.Header
	localSnapfile *LocalFileLink

	metrics     blockmetrics.Metrics
	persistence storage.StorageProvider
}

func newTemplateFromStorage(
	templateId,
	buildId,
	kernelVersion,
	firecrackerVersion string,
	memfileHeader *header.Header,
	rootfsHeader *header.Header,
	persistence storage.StorageProvider,
	metrics blockmetrics.Metrics,
	localSnapfile *LocalFileLink,
) (*storageTemplate, error) {
	files, err := storage.TemplateFiles{
		TemplateID:         templateId,
		BuildID:            buildId,
		KernelVersion:      kernelVersion,
		FirecrackerVersion: firecrackerVersion,
	}.CacheFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to create template cache files: %w", err)
	}

	return &storageTemplate{
		files:         files,
		localSnapfile: localSnapfile,
		memfileHeader: memfileHeader,
		rootfsHeader:  rootfsHeader,
		metrics:       metrics,
		persistence:   persistence,
		memfile:       utils.NewSetOnce[block.ReadonlyDevice](),
		rootfs:        utils.NewSetOnce[block.ReadonlyDevice](),
		snapfile:      utils.NewSetOnce[File](),
	}, nil
}

func (t *storageTemplate) Fetch(ctx context.Context, buildStore *build.DiffStore) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() error {
		defer wg.Done()
		if t.localSnapfile != nil {
			return t.snapfile.SetValue(t.localSnapfile)
		}

		snapfile, snapfileErr := newStorageFile(
			ctx,
			t.persistence,
			t.files.StorageSnapfilePath(),
			t.files.CacheSnapfilePath(),
		)
		if snapfileErr != nil {
			errMsg := fmt.Errorf("failed to fetch snapfile: %w", snapfileErr)

			return t.snapfile.SetError(errMsg)
		}

		return t.snapfile.SetValue(snapfile)
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()

		memfileStorage, memfileErr := NewStorage(
			ctx,
			buildStore,
			t.files.BuildID,
			build.Memfile,
			t.memfileHeader,
			t.persistence,
			t.metrics,
		)

		if memfileErr != nil {
			errMsg := fmt.Errorf("failed to create memfile storage: %w", memfileErr)

			return t.memfile.SetError(errMsg)
		}

		return t.memfile.SetValue(memfileStorage)
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()

		rootfsStorage, rootfsErr := NewStorage(
			ctx,
			buildStore,
			t.files.BuildID,
			build.Rootfs,
			t.rootfsHeader,
			t.persistence,
			t.metrics,
		)
		if rootfsErr != nil {
			errMsg := fmt.Errorf("failed to create rootfs storage: %w", rootfsErr)

			return t.rootfs.SetError(errMsg)
		}

		return t.rootfs.SetValue(rootfsStorage)
	}()

	wg.Wait()
}

func (t *storageTemplate) Close() error {
	return closeTemplate(t)
}

func (t *storageTemplate) Files() storage.TemplateCacheFiles {
	return t.files
}

func (t *storageTemplate) Memfile() (block.ReadonlyDevice, error) {
	return t.memfile.Wait()
}

func (t *storageTemplate) Rootfs() (block.ReadonlyDevice, error) {
	return t.rootfs.Wait()
}

func (t *storageTemplate) Snapfile() (File, error) {
	return t.snapfile.Wait()
}

func (t *storageTemplate) ReplaceMemfile(memfile block.ReadonlyDevice) error {
	m := utils.NewSetOnce[block.ReadonlyDevice]()
	m.SetValue(memfile)
	t.memfile = m
	return nil
}
