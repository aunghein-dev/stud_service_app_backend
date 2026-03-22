package receipt

import "context"

// Renderer is an extension point for PDF/file export integration.
// Implementations can be plugged in without changing receipt storage APIs.
type Renderer interface {
	Render(ctx context.Context, payload map[string]any) ([]byte, error)
}
