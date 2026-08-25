package cfglayer

import (
    "errors"

    "github.com/LYH2263/go-cfglayer/internal/layerstore"
    "github.com/LYH2263/go-cfglayer/internal/merge"
)

func layerToStore(l Layer) layerstore.Layer {
    vals := make(map[string]string, len(l.Values))
    for k, v := range l.Values {
        vals[k] = v
    }
    return layerstore.Layer{
        ID:       l.ID,
        Priority: l.Priority,
        Source:   l.Source,
        Values:   vals,
        Created:  l.Created,
    }
}

func layerFromStore(l layerstore.Layer) Layer {
    vals := make(map[string]string, len(l.Values))
    for k, v := range l.Values {
        vals[k] = v
    }
    return Layer{
        ID:       l.ID,
        Priority: l.Priority,
        Source:   l.Source,
        Values:   vals,
        Created:  l.Created,
    }
}

func layersToMerge(layers []Layer) []merge.Layer {
    out := make([]merge.Layer, len(layers))
    for i, l := range layers {
        vals := make(map[string]string, len(l.Values))
        for k, v := range l.Values {
            vals[k] = v
        }
        out[i] = merge.Layer{
            ID:       l.ID,
            Priority: l.Priority,
            Source:   l.Source,
            Values:   vals,
            Created:  l.Created,
        }
    }
    return out
}

func layersFromStore(layers []layerstore.Layer) []Layer {
    out := make([]Layer, len(layers))
    for i, l := range layers {
        out[i] = layerFromStore(l)
    }
    return out
}

// cloneLayers returns a fully independent copy of a layer slice: a new backing
// array with a freshly allocated Values map per element. It is used when handing
// out a view of the internal listCache so a caller cannot mutate the cached
// layer IDs / values (e.g. editing an ID as a marker) and have that edit leak
// back into the cache, corrupting the authoritative stack view.
func cloneLayers(layers []Layer) []Layer {
    out := make([]Layer, len(layers))
    for i, l := range layers {
        vals := make(map[string]string, len(l.Values))
        for k, v := range l.Values {
            vals[k] = v
        }
        out[i] = Layer{
            ID:       l.ID,
            Priority: l.Priority,
            Source:   l.Source,
            Values:   vals,
            Created:  l.Created,
        }
    }
    return out
}

func stepFromMerge(s merge.KeyStep) KeyStep {
    return KeyStep{
        LayerID:  s.LayerID,
        Priority: s.Priority,
        Value:    s.Value,
        Source:   s.Source,
    }
}

func mapStoreErr(err error) error {
    if err == nil {
        return nil
    }
    switch {
    case errors.Is(err, layerstore.ErrDuplicate):
        return ErrDuplicateLayer
    case errors.Is(err, layerstore.ErrEmpty):
        return ErrEmptyStack
    case errors.Is(err, layerstore.ErrFull):
        return ErrBadInput
    default:
        return err
    }
}
