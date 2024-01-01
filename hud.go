package main

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// HACK: To not allow dragging multiple windows at the same time
var DraggingWindow *Window

type GuiElement interface {
	Draw(offset rl.Vector2) bool
	HandleClick()
}

// Place were windows are drawn
type Hud struct {
	Visible      bool
	Windows      []GuiElement
	UpdateHandle func(*Hud)
	DrawHandle   func(*Hud)
}

func NewHud() *Hud {
	return new(Hud)
}

func (h *Hud) Add(element GuiElement) {
	h.Windows = append(h.Windows, element)
}

func (h *Hud) Remove(element GuiElement) {
	for i, e := range h.Windows {
		if e == element {
			h.Windows = append(h.Windows[:i], h.Windows[i+1:]...)
		}
	}
}

func (h *Hud) Update() {
	h.UpdateHandle(h)
}

func (h *Hud) Draw() {
	if !h.Visible {
		return
	}

	h.DrawHandle(h)

	for _, element := range h.Windows {
		if element.Draw(rl.NewVector2(0, 0)) {
			element.HandleClick()
		}
	}
}

type Window struct {
	Element
	Drag     bool
	Visible  bool
	Elements []GuiElement
	Handle   func(w *Window)
}

type Overlay struct {
	Element
	DrawHangle func(o *Overlay) bool
	Handle     func()
}

var SelectedEntity EntityInterface

func NewEntityOverlay() *Overlay {
	return &Overlay{
		Element: Element{
			Position: rl.NewVector2(0, 0),
			Size:     rl.NewVector2(WindowWidth, WindowHeight),
		},
		DrawHangle: func(o *Overlay) bool {
			for _, e := range Entities.Entities {
				rect := e.GetRectangle()
				gui.GroupBox(rect, fmt.Sprintf("%T", e))

				mousePosition := rl.GetMousePosition()
				if rl.CheckCollisionPointRec(mousePosition, e.GetRectangle()) || SelectedEntity == e {
					gui.GroupBox(rl.NewRectangle(rect.X+48, rect.Y+10, 100, 140), "Entity")
					gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+20, 96, 30), fmt.Sprintf("X: %.2f", rect.X))
					gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+60, 96, 30), fmt.Sprintf("Y: %.2f", rect.Y))
					gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+90, 96, 30), fmt.Sprintf("W: %.2f", rect.Width))
					gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+120, 96, 30), fmt.Sprintf("H: %.2f", rect.Height))

					if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
						if SelectedEntity == e {
							SelectedEntity = nil
						} else {
							SelectedEntity = e
						}
					}
				}
			}

			return false
		},
		Handle: func() {
			// TODO
		},
	}
}

func NewEmitterOverlay() *Overlay {
	return &Overlay{
		Element: Element{
			Position: rl.NewVector2(0, 0),
			Size:     rl.NewVector2(WindowWidth, WindowHeight),
		},
		DrawHangle: func(o *Overlay) bool {
			for _, e := range Emitters.Emitters {
				for _, em := range e.Emitters {
					for _, p := range em.Particles {
						rect := rl.NewRectangle(p.Position.X-em.Offset.X, p.Position.Y-em.Offset.Y, float32(em.Config.Texture.Width), float32(em.Config.Texture.Height))
						gui.GroupBox(rect, "")

						mousePosition := rl.GetMousePosition()
						if rl.CheckCollisionPointRec(mousePosition, rect) {
							gui.GroupBox(rl.NewRectangle(rect.X+48, rect.Y+10, 100, 240), "Particle")
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+20, 96, 30), fmt.Sprintf("X: %.2f", rect.X))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+60, 96, 30), fmt.Sprintf("Y: %.2f", rect.Y))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+90, 96, 30), fmt.Sprintf("W: %.2f", rect.Width))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+120, 96, 30), fmt.Sprintf("H: %.2f", rect.Height))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+150, 96, 30), fmt.Sprintf("Age: %.2f", p.Age))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+180, 96, 30), fmt.Sprintf("TTL: %.2f", p.TTL))
							gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+210, 96, 30), fmt.Sprintf("Immortal: %v", p.Immortal))
							return false
						}
						// gui.Label(rl.NewRectangle(p.Position.X, p.Position.Y, 96, 30), fmt.Sprintf("X: %.2f", p.Position.X))
						// gui.Label(rl.NewRectangle(p.Position.X, p.Position.Y+30, 96, 30), fmt.Sprintf("Y: %.2f", p.Position.Y))
					}
				}
				// rect := e.GetRectangle()
				// gui.GroupBox(rect, fmt.Sprintf("%T", e))

				// mousePosition := rl.GetMousePosition()
				// if rl.CheckCollisionPointRec(mousePosition, e.GetRectangle()) || SelectedEntity == e {
				// 	gui.GroupBox(rl.NewRectangle(rect.X+48, rect.Y+10, 100, 140), "Entity")
				// 	gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+20, 96, 30), fmt.Sprintf("X: %.2f", rect.X))
				// 	gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+60, 96, 30), fmt.Sprintf("Y: %.2f", rect.Y))
				// 	gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+90, 96, 30), fmt.Sprintf("W: %.2f", rect.Width))
				// 	gui.Label(rl.NewRectangle(rect.X+48+4, rect.Y+120, 96, 30), fmt.Sprintf("H: %.2f", rect.Height))

				// 	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				// 		if SelectedEntity == e {
				// 			SelectedEntity = nil
				// 		} else {
				// 			SelectedEntity = e
				// 		}
				// 	}
				// }
			}

			return false
		},
		Handle: func() {
			// TODO
		},
	}
}

func (o *Overlay) Draw(offset rl.Vector2) bool {
	return o.DrawHangle(o)
}

func (o *Overlay) HandleClick() {
	o.Handle()
}

func NewWindow() *Window {
	return new(Window)
}

func (w *Window) Add(element GuiElement) {
	w.Elements = append(w.Elements, element)
}

func (w *Window) Draw(offset rl.Vector2) bool {
	if !w.Visible {
		return false
	}

	dragOffset := rl.NewVector2(0, 0)
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) || w.Drag {
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(w.Position.X+offset.X, w.Position.Y+offset.Y, w.Size.X, w.Size.Y)) || w.Drag {
			if DraggingWindow == nil || DraggingWindow == w {
				w.Drag = true
				dragOffset = rl.Vector2Divide(rl.GetMouseDelta(), MouseScale)
				DraggingWindow = w
			}
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		w.Drag = false
		DraggingWindow = nil
	}

	if w.Drag {
		w.Position.X += dragOffset.X
		w.Position.Y += dragOffset.Y
	}

	if w.Position.X < 16 {
		w.Position.X = 16
	}

	if w.Position.Y < 16 {
		w.Position.Y = 16
	}

	if w.Position.X+w.Size.X > WindowWidth-16 {
		w.Position.X = WindowWidth - w.Size.X - 16
	}

	if w.Position.Y+w.Size.Y > WindowHeight-16 {
		w.Position.Y = WindowHeight - w.Size.Y - 16
	}

	rl.DrawRectangleRec(rl.NewRectangle(w.Position.X+ShadowOffset+offset.X, w.Position.Y+ShadowOffset+offset.Y, w.Size.X, w.Size.Y), rl.NewColor(0, 0, 0, 140))
	b := gui.WindowBox(rl.NewRectangle(w.Position.X+offset.X, w.Position.Y+offset.Y, w.Size.X, w.Size.Y), w.Title)

	for _, element := range w.Elements {
		if element.Draw(w.Position) {
			element.HandleClick()
		}
	}

	return b
}

func (w *Window) HandleClick() {
	w.Handle(w)
}

type Element struct {
	Position rl.Vector2
	Size     rl.Vector2
	Title    string
}

type Button struct {
	Element
	Shadow bool
	Handle func()
}

func (b *Button) Draw(offset rl.Vector2) bool {
	if b.Shadow {
		rl.DrawRectangleRec(rl.NewRectangle(b.Position.X+ShadowOffset+offset.X, b.Position.Y+ShadowOffset+offset.Y, b.Size.X, b.Size.Y), rl.NewColor(0, 0, 0, 140))
	}

	return gui.Button(rl.NewRectangle(b.Position.X+offset.X, b.Position.Y+offset.Y, b.Size.X, b.Size.Y), b.Title)
}

func (b *Button) HandleClick() {
	b.Handle()
}

type EmittersList struct {
	Element
}

func (e *EmittersList) Draw(offset rl.Vector2) bool {
	emitters := Emitters.Emitters
	for i, emitter := range emitters {
		if gui.Button(rl.NewRectangle(e.Position.X+offset.X, e.Position.Y+offset.Y+(float32(i)*32), e.Size.X, 32), fmt.Sprintf("%T", emitter)) {
			Debug.Add(NewDebugEmitterWindow(rl.NewVector2((WindowWidth/2)-100, (WindowHeight/2)-100), rl.NewVector2(200, 200), emitter))
		}
	}

	return false
}

func (e *EmittersList) HandleClick() {
}

type ToggleBox struct {
	Element
	Handle func()
}

func (t *ToggleBox) Draw(offset rl.Vector2) bool {
	return gui.Button(rl.NewRectangle(t.Position.X+offset.X, t.Position.Y+offset.Y, t.Size.X, t.Size.Y), t.Title)
}

func (t *ToggleBox) HandleClick() {
	t.Handle()
}

type IntValue struct {
	Element
	Value    func() int
	Handle   func()
	MinValue int
	MaxValue int
	Edit     bool
}

func (i *IntValue) Draw(offset rl.Vector2) bool {
	v := int32(i.Value())
	return gui.ValueBox(rl.NewRectangle(i.Position.X+offset.X, i.Position.Y+offset.Y, i.Size.X, i.Size.Y), i.Title, &v, i.MinValue, i.MaxValue, i.Edit)
}

func (i *IntValue) HandleClick() {
	i.Handle()
}

type BoolValue struct {
	Element
	Value  func() bool
	Handle func()
}

func (i *BoolValue) Draw(offset rl.Vector2) bool {
	return gui.CheckBox(rl.NewRectangle(i.Position.X+offset.X, i.Position.Y+offset.Y, i.Size.X, i.Size.Y), i.Title, i.Value())
}

func (i *BoolValue) HandleClick() {
	i.Handle()
}

const ShadowOffset = 8

func NewDebugHud() *Hud {
	return &Hud{
		Windows: []GuiElement{
			NewEntityOverlay(),
		},
		UpdateHandle: func(h *Hud) {
			if rl.IsKeyPressed(rl.KeyF1) {
				h.Visible = !h.Visible
			}

			if h.Visible && rl.IsKeyPressed(rl.KeyEnter) {
				Debug.Add(NewDebugWindow(rl.NewVector2((WindowWidth/2)-100, (WindowHeight/2)-100), rl.NewVector2(200, 200)))
			}
		},
		DrawHandle: func(h *Hud) {
			rl.DrawTextEx(*Assets.FontsManager.SmallFont, "Debug", rl.NewVector2(10, 10), SmallFontSize, 0, rl.RayWhite)
			rl.DrawTextEx(*Assets.FontsManager.SmallFont, fmt.Sprintf("Frames: %d", rl.GetFPS()), rl.NewVector2(WindowWidth-130, 10), SmallFontSize, 0, rl.RayWhite)
		},
	}
}

func NewDebugWindow(position rl.Vector2, size rl.Vector2) *Window {
	return &Window{
		Element: Element{
			Position: position,
			Size:     size,
			Title:    "Debug",
		},
		Visible: true,
		Elements: []GuiElement{
			&Button{
				Element: Element{
					Position: rl.NewVector2(8, 32),
					Size:     rl.NewVector2(184, 30),
					Title:    "Emitter",
				},
				Handle: func() {
					Debug.Add(NewDebugEmittersWindow(position, size))
				},
			},
		},
		Handle: func(w *Window) {
			w.Visible = !w.Visible
		},
	}
}

func NewDebugEmittersWindow(position rl.Vector2, size rl.Vector2) *Window {
	Debug.Add(NewEmitterOverlay())
	return &Window{
		Element: Element{
			Position: position,
			Size:     size,
			Title:    "Emitters Manager",
		},
		Visible: true,
		Elements: []GuiElement{
			&Button{
				Element: Element{
					Position: rl.NewVector2(8, 32),
					Size:     rl.NewVector2(184, 30),
					Title:    "Simulate",
				},
				Handle: func() {
					Emitters.Updates = !Emitters.Updates
				},
			},
			&IntValue{
				Element: Element{
					Position: rl.NewVector2(48, 68),
					Size:     rl.NewVector2(184-40, 30),
					Title:    "Count ",
				},
				Value: func() int {
					return Emitters.Count()
				},
				Handle: func() {

				},
				MinValue: 0,
				MaxValue: 10000,
				Edit:     false,
			},
			&BoolValue{
				Element: Element{
					Position: rl.NewVector2(8, 108),
					Size:     rl.NewVector2(184-40, 30),
					Title:    "Active  ",
				},
				Value: func() bool {
					return Emitters.Updates
				},
				Handle: func() {
				},
			},
			&EmittersList{
				Element: Element{
					Position: rl.NewVector2(8, 148),
					Size:     rl.NewVector2(184-40, 30),
				},
			},
		},
		Handle: func(w *Window) {
			Debug.Remove(w)
		},
	}
}

func NewDebugEmitterWindow(position rl.Vector2, size rl.Vector2, e *ParticleSystem) *Window {
	return &Window{
		Element: Element{
			Position: position,
			Size:     size,
			Title:    "Emitter",
		},
		Visible: true,
		Handle: func(w *Window) {
			w.Visible = !w.Visible
		},
		Elements: []GuiElement{
			&Button{
				Element: Element{
					Position: rl.NewVector2(8, 32),
					Size:     rl.NewVector2(184, 30),
					Title:    "Remove",
				},
				Handle: func() {
					Emitters.Remove(e)
				},
			},
		},
	}
}
