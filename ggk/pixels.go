package ggk

import (
	"log"
	"sync"
)

func validatePixelsColorTable(info ImageInfo, ct *ColorTable) {
	if info.IsEmpty() {
		return // can't require ct if the dimensions are empty.
	}

	if info.ColorType() == KColorTypeIndex8 {
		if ct == nil {
			log.Printf(`WARNING: validatePixelsColorTable ct is nil for ColorTypeIndex8`)
		}
	} else {
		if ct != nil {
			log.Printf(`WARNING: validatePixelsColorTable ct is not nil for non-ColorTypeIndex8`)
		}
	}
}

type PixelsVirt interface {
	// OnNewLockPixels returns true and fills out the LockRec for the pixels on
	// success, returns false and ignores the LockRec parameter on failure.
	//
	// The caller will have already acquired a mutex for thread safety, so this
	// method need not do that.
	OnNewLockPixels(*PixelsLockRec) bool

	// OnUnlockPixels balancing the previous successful call to OnNewLockPixels.
	// Tht locked pixel address will no longer be referenced, so the subclass is
	// free to move or discard that memory.
	//
	// The caller will have already acquired a mutex for thread safety, so this
	// method need not do that.
	OnUnlockPixels()

	// OnUnlockPixels default impl returns true.
	OnLockPixelsAreWritable() bool

	// OnReadPixels for pixelRefs that don't have access to their raw pixels,
	// they may be able to make a copy of them (e.g. if the pixels are on the
	// GPU).
	//
	// The base class implementation returns false.
	OnReadPixels(dst *Bitmap, subsetOrNil *Rect) bool

	// OnRefEncodedData default impl returns nil.
	// OnRefEncodedData() *Data

	// OnNotifyPixelsChanged default impl does nothing.
	OnNotifyPixelsChanged()

	// OnGetYUV8Planes default impl returns false.
	// OnGetYUV8Planes(size [3]Size, planes [3]interface{}, rowBytes [3]int, cs YUVColorSpace) bool

	// GetAllocatedSizeInBytes returns the size (in bytes) of the internally
	// allocated memory. This should be implemented in all serializable PixelRef
	// derived classes. Bitmap::pixelRefOffset + Bitmap::GetSafeSize() should
	// never overflow this value, otherwise the rendering code may attempt to
	// read memory out of bounds.
	GetAllocatedSizeInBytes() uint

	// OnRequestLock send request for the pixels access lock.
	// OnRequestLock(*LockRequest, *LockResult) bool
}

type Pixels struct {
	// Virt is virtual interface, set Virt in NewXYZ() call Pixels.Virt.Func()
	// to get the virtual aliblity.
	Virt PixelsVirt

	prelocked bool
	mutex     sync.Mutex
	lockRec   PixelsLockRec
	lockCount int
}

func NewPixels() *Pixels {
	// TOIMPL
	var p = new(Pixels)
	return p
}

// Pixels return the pixel memory bytes returned from LockPixels, or nil if the
// lockCount is 0.
func (p *Pixels) Pixels() []byte {
	return p.lockRec.pixels
}

// ColorTable return the current colorTable (if any) if pixels are locked, or
// nil.
func (p *Pixels) ColorTable() *ColorTable {
	return p.lockRec.colorTable
}

// RowBytes return the current rowBytes (if any) if pixels are locked, or nil.
func (p *Pixels) RowBytes() int {
	return p.lockRec.rowBytes
}

// Just need a > 0 value, so pick a funny one to aid in debugging.
const kPixelsPrelockedLockCount = 123456789

// LockPixels try to get the pixels lock, and prepare for read/write the pixels.
// For the historical reasons, we always inc lockCount, even if we return false.
// It would be nice to change this (it seems), and only inc if we actually
// succeeding...
func (p *Pixels) LockPixels() bool {
	if p.prelocked && p.lockCount != kPixelsPrelockedLockCount {
		log.Printf(`WARNING: Pixels.LockPixels prelocked and lockCount is not matched.`)
	}

	if !p.prelocked {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		var ok = p.LockPixelsInsideMutex()
		if !ok {
			p.lockCount += 1
			return false
		}
	}

	return false
}

func (p *Pixels) LockPixelsToRec() (bool, *PixelsLockRec) {
	if p.LockPixels() {
		return true, &(p.lockRec)
	}
	return false, nil
}

// Increments lockCount only on success.
func (p *Pixels) LockPixelsInsideMutex() bool {
	p.lockCount++

	if p.lockCount == 1 {
		if !p.Virt.OnNewLockPixels(&p.lockRec) {
			p.lockRec.SetZero()
			p.lockCount -= 1 // We return lockCount unchanged if we faile.
			return false
		}
	}

	return false
}

func (p *Pixels) UnlockPixels() {
	if !p.prelocked {
		p.mutex.Lock()
		defer p.mutex.Unlock()

		p.lockCount--
		if p.lockCount == 0 {
			// Don't call OnUnlockPixels unless OnLockPixels succeeded.
			if p.lockRec.pixels != nil {
				p.Virt.OnUnlockPixels()
				p.lockRec.SetZero()
			}
		}
	}
}

func (p *Pixels) LockPixelsAreWritable() bool {
	return p.Virt.OnLockPixelsAreWritable()
}

// OnLockPixelsAreWritable default impl return true.
func (p *Pixels) OnLockPixelsAreWritable() bool {
	return true
}

func (p *Pixels) ReadPixels(dst *Bitmap, subset *Rect) bool {
	return p.Virt.OnReadPixels(dst, subset)
}

func (p *Pixels) OnReadPixels(dst *Bitmap, subset *Rect) bool {
	return false
}

func (p *Pixels) OnNotifyPixelsChanged() {
	// empty
}

// func (p *Pixels) OnRefEncodedData() *Data {
// 	return nil
// }

func (p *Pixels) GetAllocatedSizeInBytes() uint {
	return 0
}

func (p *Pixels) SetPrelocked(pixels []byte) {
	// only call me in your constructor, other wise fLockCount tracking can get
	// out of sync.
	p.lockRec.pixels = pixels
	p.lockCount = kPixelsPrelockedLockCount
	p.prelocked = true
}

// PixelRefLockRec to access the actual pixels of a pixelRef, it must be
// "locked". Calling LockPixels returns a PixelRefLockRec struct (on success).
type PixelsLockRec struct {
	pixels     []byte
	colorTable *ColorTable
	rowBytes   int
}

func (r *PixelsLockRec) SetZero() {
	var zero PixelsLockRec
	*r = zero
}

func (r PixelsLockRec) IsZero() bool {
	return r.pixels == nil && r.colorTable == nil && r.rowBytes == 0
}

type MemPixelsFactory struct {
}

func MemPixelsDefaultFactory() *MemPixelsFactory {
	// TOIMPL
	return nil
}

func (f *MemPixelsFactory) Create(imageInfo ImageInfo, rowBytes int, colorTable *ColorTable) *Pixels {
	// TOIMPL
	return nil
}

type MemPixels struct {
	*Pixels
	storage []byte
}

// NewMemPixels return a new instance with the provided stroage, rowBytes,
// and optional colortable. The caller is responsible for managing the
// lifetime of the pixel storage buffer, as this pixels will not try
// to delete it.
//
// The pixels will ref the colortable (if not nil)
//
// Returns nil on failture.
func NewMemPixels(storage []byte) *MemPixels {
	var p MemPixels
	p.Pixels = NewPixels()
	p.storage = storage
	p.SetPrelocked(storage)
	return &p
}

func NewMemPixelsAlloc(info ImageInfo, rowBytes int) *MemPixels {
	var p MemPixels
	p.Pixels = NewPixels()
	p.storage = make([]byte, info.SafeSize(rowBytes))
	p.SetPrelocked(p.storage)
	return &p
}

func isImageInfoValid(info ImageInfo) bool {
	if info.Width() < 0 || info.Height() < 0 || !info.ColorType().IsVaild() ||
		!info.AlphaType().IsValid() {
		return false
	}

	return true
}
