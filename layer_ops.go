package cfglayer

import (
	"context"
	"time"

	"github.com/LYH2263/go-cfglayer/internal/validate"
)

func (m *Merger) PushLayer(ctx context.Context, layer Layer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := m.checkOpen(); err != nil {
		return err
	}
	if err := validate.Layer(validate.LayerInput{ID: layer.ID, Values: layer.Values}); err != nil {
		return ErrBadInput
	}
	if layer.Values == nil {
		layer.Values = map[string]string{}
	}
	vals := make(map[string]string, len(layer.Values))
	for k, v := range layer.Values {
		vals[k] = v
	}
	layer.Values = vals
	if layer.Created.IsZero() {
		layer.Created = time.Now()
	}
	if err := mapStoreErr(m.store.Push(layerToStore(layer))); err != nil {
		return err
	}
	m.mu.Lock()
	m.listCache = nil
	m.mu.Unlock()
	if m.audit != nil {
		m.audit.Printf("push layer=%s keys=%d", layer.ID, len(layer.Values))
	}
	return nil
}

func (m *Merger) PopLayer() (Layer, error) {
	if err := m.checkOpen(); err != nil {
		return Layer{}, err
	}
	layer, err := m.store.Pop()
	if err != nil {
		return Layer{}, mapStoreErr(err)
	}
	m.mu.Lock()
	m.listCache = nil
	m.mu.Unlock()
	if m.audit != nil {
		m.audit.Printf("pop layer=%s", layer.ID)
	}
	return layerFromStore(layer), nil
}

func (m *Merger) ListLayers() []Layer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.listCache == nil {
		m.listCache = layersFromStore(m.store.List())
	}
	// Never hand out the internal cache slice: a caller mutating the
	// returned elements (e.g. editing an ID) would otherwise leak the
	// change back into m.listCache, corrupting the real stack view.
	return cloneLayers(m.listCache)
}
