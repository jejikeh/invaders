package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const ShadowOffset = 8

type DebugHud struct {
	Window        *FloatWindow
	Visible       bool
	EmitterButton Button
	EmitterWindow *EmitterWindow
}

func NewDebugHud() *DebugHud {
	return &DebugHud{
		Window: &FloatWindow{
			Position: rl.NewVector2(WindowWidth/2-100, WindowHeight/2-100),
			Size:     rl.NewVector2(200, 200),
			Visible:  true,
			Titile:   "Debug",
		},
		EmitterButton: Button{
			Position: rl.NewVector2(0, 0),
			Size:     rl.NewVector2(184, 30),
			Title:    "Emitter",
			Shadow:   false,
		},
		Visible:       true,
		EmitterWindow: NewEmitterWindow(),
	}
}

type Button struct {
	Position rl.Vector2
	Size     rl.Vector2
	Title    string
	Shadow   bool
}

func (d *DebugHud) FirstDraw() {

}

func (b *Button) Draw() bool {
	if b.Shadow {
		rl.DrawRectangleRec(rl.NewRectangle(b.Position.X+ShadowOffset, b.Position.Y+ShadowOffset, b.Size.X, b.Size.Y), rl.NewColor(0, 0, 0, 140))
	}
	return gui.Button(rl.NewRectangle(b.Position.X, b.Position.Y, b.Size.X, b.Size.Y), b.Title)
}

type FloatWindow struct {
	Position   rl.Vector2
	Size       rl.Vector2
	Drag       bool
	DragOffset rl.Vector2
	Visible    bool
	Titile     string
}

type EmitterWindow struct {
	FloatWindow
}

func NewEmitterWindow() *EmitterWindow {
	return &EmitterWindow{
		FloatWindow: FloatWindow{
			Position:   rl.NewVector2(0, 0),
			Size:       rl.NewVector2(200, 200),
			Drag:       false,
			DragOffset: rl.NewVector2(0, 0),
			Visible:    true,
			Titile:     "Emitter",
		},
	}
}

func (e *EmitterWindow) Draw() bool {
	if !e.Visible {
		return false
	}

	shouldClose := e.FloatWindow.Draw()

	value := int32(Emitters.Count())
	gui.ValueBox(rl.NewRectangle(e.Position.X+48, e.Position.Y+32, e.Size.X-72, 24), "Count: ", &value, 0, 100000, false)

	return shouldClose
}

func (e *FloatWindow) Draw() bool {
	if !e.Visible {
		return false
	}

	dragOffset := rl.NewVector2(0, 0)
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(e.Position.X, e.Position.Y, e.Size.X, e.Size.Y)) || e.Drag {
			e.Drag = true
			dragOffset = rl.Vector2Divide(rl.GetMouseDelta(), rl.NewVector2(1, 1))
			// dragOffset = rl.Vector2Subtract(rl.GetMouseDelta(), rl.NewVector2(MouseScale, MouseScale))
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		e.Drag = false
	}

	if e.Drag {
		e.Position.X += dragOffset.X
		e.Position.Y += dragOffset.Y
	}

	if e.Position.X < 8 {
		e.Position.X = 8
	}

	if e.Position.Y < 8 {
		e.Position.Y = 8
	}

	if e.Position.X+e.Size.X > float32(rl.GetScreenWidth())-8 {
		e.Position.X = float32(rl.GetScreenWidth()) - e.Size.X - 8
	}

	if e.Position.Y+e.Size.Y > float32(rl.GetScreenHeight())-8 {
		e.Position.Y = float32(rl.GetScreenHeight()) - e.Size.Y - 8
	}

	rl.DrawRectangleRec(rl.NewRectangle(e.Position.X+ShadowOffset, e.Position.Y+ShadowOffset, e.Size.X, e.Size.Y), rl.NewColor(0, 0, 0, 140))
	if gui.WindowBox(rl.NewRectangle(e.Position.X, e.Position.Y, e.Size.X, e.Size.Y), e.Titile) {
		e.Visible = false
	}

	return e.Visible
}

func (d *DebugHud) Start() {
	d.Visible = false
}

func (d *DebugHud) Draw() {
	if !d.Visible {
		return
	}

	if d.Visible {
		d.Window.Visible = true
	}

	if d.Window.Draw() {
		d.EmitterButton.Position = rl.NewVector2(d.Window.Position.X+8, d.Window.Position.Y+32)
		d.EmitterButton.Draw()
	} else {
		d.Visible = false
	}

	if d.EmitterWindow != nil {
		d.EmitterWindow.Draw()
	}
}

func (d *DebugHud) Update() {
	if rl.IsKeyPressed(rl.KeyF1) {
		d.Visible = !d.Visible
	}
}

func (d *DebugHud) GetLayer() int {
	return 100
}
