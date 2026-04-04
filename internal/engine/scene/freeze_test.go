package scene

import "testing"

func TestFreezeController(t *testing.T) {
	tests := []struct {
		name  string
		setup func(f *FreezeController)
		want  bool
	}{
		{
			name:  "no freeze",
			setup: func(_ *FreezeController) {},
			want:  false,
		},
		{
			name:  "freeze(3) frame 0 is frozen",
			setup: func(f *FreezeController) { f.FreezeFrame(3) },
			want:  true,
		},
		{
			name:  "freeze(3) frame 1 is frozen",
			setup: func(f *FreezeController) { f.FreezeFrame(3); f.Tick() },
			want:  true,
		},
		{
			name:  "freeze(3) frame 2 is frozen",
			setup: func(f *FreezeController) { f.FreezeFrame(3); f.Tick(); f.Tick() },
			want:  true,
		},
		{
			name:  "freeze(3) frame 3 is not frozen",
			setup: func(f *FreezeController) { f.FreezeFrame(3); f.Tick(); f.Tick(); f.Tick() },
			want:  false,
		},
		{
			name:  "freeze(0) no-op",
			setup: func(f *FreezeController) { f.FreezeFrame(0) },
			want:  false,
		},
		{
			name:  "freeze(-1) no-op",
			setup: func(f *FreezeController) { f.FreezeFrame(-1) },
			want:  false,
		},
		{
			name: "reset mid-freeze latest call wins",
			setup: func(f *FreezeController) {
				f.FreezeFrame(3)
				f.Tick()
				f.FreezeFrame(5)
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FreezeController{}
			tt.setup(f)
			if got := f.IsFrozen(); got != tt.want {
				t.Errorf("IsFrozen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFreezeControllerResetWins(t *testing.T) {
	f := &FreezeController{}
	f.FreezeFrame(2)
	f.Tick()
	f.FreezeFrame(4)

	for i := 0; i < 4; i++ {
		if !f.IsFrozen() {
			t.Fatalf("expected frozen at tick %d", i)
		}
		f.Tick()
	}

	if f.IsFrozen() {
		t.Error("expected not frozen after 4 ticks following FreezeFrame(4)")
	}
}
