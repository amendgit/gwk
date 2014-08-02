// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysc

// MSG.
type MSG struct {
	HWnd    Handle // HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// TRUE and FALSE
const (
	FALSE = 0
	TRUE  = 1
)

// Window Message.
const (
	WM_APP                    = 32768
	WM_ACTIVATE               = 6
	WM_ACTIVATEAPP            = 28
	WM_AFXFIRST               = 864
	WM_AFXLAST                = 895
	WM_ASKCBFORMATNAME        = 780
	WM_CANCELJOURNAL          = 75
	WM_CANCELMODE             = 31
	WM_CAPTURECHANGED         = 533
	WM_CHANGECBCHAIN          = 781
	WM_CHAR                   = 258
	WM_CHARTOITEM             = 47
	WM_CHILDACTIVATE          = 34
	WM_CLEAR                  = 771
	WM_CLOSE                  = 16
	WM_COMMAND                = 273
	WM_COMMNOTIFY             = 68 /* OBSOLETE */
	WM_COMPACTING             = 65
	WM_COMPAREITEM            = 57
	WM_CONTEXTMENU            = 123
	WM_COPY                   = 769
	WM_COPYDATA               = 74
	WM_CREATE                 = 1
	WM_CTLCOLORBTN            = 309
	WM_CTLCOLORDLG            = 310
	WM_CTLCOLOREDIT           = 307
	WM_CTLCOLORLISTBOX        = 308
	WM_CTLCOLORMSGBOX         = 306
	WM_CTLCOLORSCROLLBAR      = 311
	WM_CTLCOLORSTATIC         = 312
	WM_CUT                    = 768
	WM_DEADCHAR               = 259
	WM_DELETEITEM             = 45
	WM_DESTROY                = 2
	WM_DESTROYCLIPBOARD       = 775
	WM_DEVICECHANGE           = 537
	WM_DEVMODECHANGE          = 27
	WM_DISPLAYCHANGE          = 126
	WM_DRAWCLIPBOARD          = 776
	WM_DRAWITEM               = 43
	WM_DROPFILES              = 563
	WM_ENABLE                 = 10
	WM_ENDSESSION             = 22
	WM_ENTERIDLE              = 289
	WM_ENTERMENULOOP          = 529
	WM_ENTERSIZEMOVE          = 561
	WM_ERASEBKGND             = 20
	WM_EXITMENULOOP           = 530
	WM_EXITSIZEMOVE           = 562
	WM_FONTCHANGE             = 29
	WM_GETDLGCODE             = 135
	WM_GETFONT                = 49
	WM_GETHOTKEY              = 51
	WM_GETICON                = 127
	WM_GETMINMAXINFO          = 36
	WM_GETTEXT                = 13
	WM_GETTEXTLENGTH          = 14
	WM_HANDHELDFIRST          = 856
	WM_HANDHELDLAST           = 863
	WM_HELP                   = 83
	WM_HOTKEY                 = 786
	WM_HSCROLL                = 276
	WM_HSCROLLCLIPBOARD       = 782
	WM_ICONERASEBKGND         = 39
	WM_INITDIALOG             = 272
	WM_INITMENU               = 278
	WM_INITMENUPOPUP          = 279
	WM_INPUT                  = 0X00FF
	WM_INPUTLANGCHANGE        = 81
	WM_INPUTLANGCHANGEREQUEST = 80
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_KILLFOCUS              = 8
	WM_MDIACTIVATE            = 546
	WM_MDICASCADE             = 551
	WM_MDICREATE              = 544
	WM_MDIDESTROY             = 545
	WM_MDIGETACTIVE           = 553
	WM_MDIICONARRANGE         = 552
	WM_MDIMAXIMIZE            = 549
	WM_MDINEXT                = 548
	WM_MDIREFRESHMENU         = 564
	WM_MDIRESTORE             = 547
	WM_MDISETMENU             = 560
	WM_MDITILE                = 550
	WM_MEASUREITEM            = 44
	WM_GETOBJECT              = 0X003D
	WM_CHANGEUISTATE          = 0X0127
	WM_UPDATEUISTATE          = 0X0128
	WM_QUERYUISTATE           = 0X0129
	WM_UNINITMENUPOPUP        = 0X0125
	WM_MENURBUTTONUP          = 290
	WM_MENUCOMMAND            = 0X0126
	WM_MENUGETOBJECT          = 0X0124
	WM_MENUDRAG               = 0X0123
	WM_APPCOMMAND             = 0X0319
	WM_MENUCHAR               = 288
	WM_MENUSELECT             = 287
	WM_MOVE                   = 3
	WM_MOVING                 = 534
	WM_NCACTIVATE             = 134
	WM_NCCALCSIZE             = 131
	WM_NCCREATE               = 129
	WM_NCDESTROY              = 130
	WM_NCHITTEST              = 132
	WM_NCLBUTTONDBLCLK        = 163
	WM_NCLBUTTONDOWN          = 161
	WM_NCLBUTTONUP            = 162
	WM_NCMBUTTONDBLCLK        = 169
	WM_NCMBUTTONDOWN          = 167
	WM_NCMBUTTONUP            = 168
	WM_NCXBUTTONDOWN          = 171
	WM_NCXBUTTONUP            = 172
	WM_NCXBUTTONDBLCLK        = 173
	WM_NCMOUSEHOVER           = 0X02A0
	WM_NCMOUSELEAVE           = 0X02A2
	WM_NCMOUSEMOVE            = 160
	WM_NCPAINT                = 133
	WM_NCRBUTTONDBLCLK        = 166
	WM_NCRBUTTONDOWN          = 164
	WM_NCRBUTTONUP            = 165
	WM_NEXTDLGCTL             = 40
	WM_NEXTMENU               = 531
	WM_NOTIFY                 = 78
	WM_NOTIFYFORMAT           = 85
	WM_NULL                   = 0
	WM_PAINT                  = 15
	WM_PAINTCLIPBOARD         = 777
	WM_PAINTICON              = 38
	WM_PALETTECHANGED         = 785
	WM_PALETTEISCHANGING      = 784
	WM_PARENTNOTIFY           = 528
	WM_PASTE                  = 770
	WM_PENWINFIRST            = 896
	WM_PENWINLAST             = 911
	WM_POWER                  = 72
	WM_POWERBROADCAST         = 536
	WM_PRINT                  = 791
	WM_PRINTCLIENT            = 792
	WM_QUERYDRAGICON          = 55
	WM_QUERYENDSESSION        = 17
	WM_QUERYNEWPALETTE        = 783
	WM_QUERYOPEN              = 19
	WM_QUEUESYNC              = 35
	WM_QUIT                   = 18
	WM_RENDERALLFORMATS       = 774
	WM_RENDERFORMAT           = 773
	WM_SETCURSOR              = 32
	WM_SETFOCUS               = 7
	WM_SETFONT                = 48
	WM_SETHOTKEY              = 50
	WM_SETICON                = 128
	WM_SETREDRAW              = 11
	WM_SETTEXT                = 12
	WM_SETTINGCHANGE          = 26
	WM_SHOWWINDOW             = 24
	WM_SIZE                   = 5
	WM_SIZECLIPBOARD          = 779
	WM_SIZING                 = 532
	WM_SPOOLERSTATUS          = 42
	WM_STYLECHANGED           = 125
	WM_STYLECHANGING          = 124
	WM_SYSCHAR                = 262
	WM_SYSCOLORCHANGE         = 21
	WM_SYSCOMMAND             = 274
	WM_SYSDEADCHAR            = 263
	WM_SYSKEYDOWN             = 260
	WM_SYSKEYUP               = 261
	WM_TCARD                  = 82
	WM_THEMECHANGED           = 794
	WM_TIMECHANGE             = 30
	WM_TIMER                  = 275
	WM_UNDO                   = 772
	WM_USER                   = 1024
	WM_USERCHANGED            = 84
	WM_VKEYTOITEM             = 46
	WM_VSCROLL                = 277
	WM_VSCROLLCLIPBOARD       = 778
	WM_WINDOWPOSCHANGED       = 71
	WM_WINDOWPOSCHANGING      = 70
	WM_WININICHANGE           = 26
	WM_KEYFIRST               = 256
	WM_KEYLAST                = 264
	WM_SYNCPAINT              = 136
	WM_MOUSEACTIVATE          = 33
	WM_MOUSEMOVE              = 512
	WM_LBUTTONDOWN            = 513
	WM_LBUTTONUP              = 514
	WM_LBUTTONDBLCLK          = 515
	WM_RBUTTONDOWN            = 516
	WM_RBUTTONUP              = 517
	WM_RBUTTONDBLCLK          = 518
	WM_MBUTTONDOWN            = 519
	WM_MBUTTONUP              = 520
	WM_MBUTTONDBLCLK          = 521
	WM_MOUSEWHEEL             = 522
	WM_MOUSEFIRST             = 512
	WM_XBUTTONDOWN            = 523
	WM_XBUTTONUP              = 524
	WM_XBUTTONDBLCLK          = 525
	WM_MOUSELAST              = 525
	WM_MOUSEHOVER             = 0X2A1
	WM_MOUSELEAVE             = 0X2A3
)

// Predifined brushes.
const (
	COLOR_3DDKSHADOW              = 21
	COLOR_3DFACE                  = 15
	COLOR_3DHILIGHT               = 20
	COLOR_3DHIGHLIGHT             = 20
	COLOR_3DLIGHT                 = 22
	COLOR_BTNHILIGHT              = 20
	COLOR_3DSHADOW                = 16
	COLOR_ACTIVEBORDER            = 10
	COLOR_ACTIVECAPTION           = 2
	COLOR_APPWORKSPACE            = 12
	COLOR_BACKGROUND              = 1
	COLOR_DESKTOP                 = 1
	COLOR_BTNFACE                 = 15
	COLOR_BTNHIGHLIGHT            = 20
	COLOR_BTNSHADOW               = 16
	COLOR_BTNTEXT                 = 18
	COLOR_CAPTIONTEXT             = 9
	COLOR_GRAYTEXT                = 17
	COLOR_HIGHLIGHT               = 13
	COLOR_HIGHLIGHTTEXT           = 14
	COLOR_INACTIVEBORDER          = 11
	COLOR_INACTIVECAPTION         = 3
	COLOR_INACTIVECAPTIONTEXT     = 19
	COLOR_INFOBK                  = 24
	COLOR_INFOTEXT                = 23
	COLOR_MENU                    = 4
	COLOR_MENUTEXT                = 7
	COLOR_SCROLLBAR               = 0
	COLOR_WINDOW                  = 5
	COLOR_WINDOWFRAME             = 6
	COLOR_WINDOWTEXT              = 8
	COLOR_HOTLIGHT                = 26
	COLOR_GRADIENTACTIVECAPTION   = 27
	COLOR_GRADIENTINACTIVECAPTION = 28
)

// Window class styles
const (
	CS_VREDRAW         = 0x00000001
	CS_HREDRAW         = 0x00000002
	CS_KEYCVTWINDOW    = 0x00000004
	CS_DBLCLKS         = 0x00000008
	CS_OWNDC           = 0x00000020
	CS_CLASSDC         = 0x00000040
	CS_PARENTDC        = 0x00000080
	CS_NOKEYCVT        = 0x00000100
	CS_NOCLOSE         = 0x00000200
	CS_SAVEBITS        = 0x00000800
	CS_BYTEALIGNCLIENT = 0x00001000
	CS_BYTEALIGNWINDOW = 0x00002000
	CS_GLOBALCLASS     = 0x00004000
	CS_IME             = 0x00010000
	CS_DROPSHADOW      = 0x00020000
)

type WNDCLASSEX struct {
	Size          uint32
	Style         uint32
	FnWndProc     uintptr
	ClassExtra    int32
	WindowExtra   int32
	HInstance     Handle
	HIcon         Handle
	HCursor       Handle
	HbrBackground Handle
	MenuName      *uint16
	ClassName     *uint16
	HIconSm       Handle
}

const CW_USEDEFAULT = 0

// Window Style Constants.
const (
	WS_BORDER              = 0x800000
	WS_CAPTION             = 0xc00000
	WS_CHILD               = 0x40000000
	WS_CHILDWINDOW         = 0x40000000
	WS_CLIPCHILDREN        = 0x2000000
	WS_CLIPSIBLINGS        = 0x4000000
	WS_DISABLED            = 0x8000000
	WS_DLGFRAME            = 0x400000
	WS_GROUP               = 0x20000
	WS_HSCROLL             = 0x100000
	WS_ICONIC              = 0x20000000
	WS_MAXIMIZE            = 0x1000000
	WS_MAXIMIZEBOX         = 0x10000
	WS_MINIMIZE            = 0x20000000
	WS_MINIMIZEBOX         = 0x20000
	WS_OVERLAPPED          = 0
	WS_OVERLAPPEDWINDOW    = 0xcf0000
	WS_POPUP               = 0x80000000
	WS_POPUPWINDOW         = 0x80880000
	WS_SIZEBOX             = 0x40000
	WS_SYSMENU             = 0x80000
	WS_TABSTOP             = 0x10000
	WS_THICKFRAME          = 0x40000
	WS_TILED               = 0
	WS_TILEDWINDOW         = 0xcf0000
	WS_VISIBLE             = 0x10000000
	WS_VSCROLL             = 0x200000
	WS_EX_ACCEPTFILES      = 16
	WS_EX_APPWINDOW        = 0x40000
	WS_EX_CLIENTEDGE       = 512
	WS_EX_COMPOSITED       = 0x2000000 /* XP */
	WS_EX_CONTEXTHELP      = 0x400
	WS_EX_CONTROLPARENT    = 0x10000
	WS_EX_DLGMODALFRAME    = 1
	WS_EX_LAYERED          = 0x80000  /* w2k */
	WS_EX_LAYOUTRTL        = 0x400000 /* w98, w2k */
	WS_EX_LEFT             = 0
	WS_EX_LEFTSCROLLBAR    = 0x4000
	WS_EX_LTRREADING       = 0
	WS_EX_MDICHILD         = 64
	WS_EX_NOACTIVATE       = 0x8000000 /* w2k */
	WS_EX_NOINHERITLAYOUT  = 0x100000  /* w2k */
	WS_EX_NOPARENTNOTIFY   = 4
	WS_EX_OVERLAPPEDWINDOW = 0x300
	WS_EX_PALETTEWINDOW    = 0x188
	WS_EX_RIGHT            = 0x1000
	WS_EX_RIGHTSCROLLBAR   = 0
	WS_EX_RTLREADING       = 0x2000
	WS_EX_STATICEDGE       = 0x20000
	WS_EX_TOOLWINDOW       = 128
	WS_EX_TOPMOST          = 8
	WS_EX_TRANSPARENT      = 32
	WS_EX_WINDOWEDGE       = 256
)

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

// ShowWindow constants.
const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
	SW_MAX             = 11
)

const (
	IDI_APPLICATION = 32512
	IDI_HAND        = 32513
	IDI_QUESTION    = 32514
	IDI_EXCLAMATION = 32515
	IDI_ASTERISK    = 32516
	IDI_WINLOGO     = 32517
)

// PeekMessage wRemoveMsg value
const (
	PM_NOREMOVE = 0x000
	PM_REMOVE   = 0x001
	PM_NOYIELD  = 0x002
)

const (
	IDC_ARROW       = 32512
	IDC_IBEAM       = 32513
	IDC_WAIT        = 32514
	IDC_CROSS       = 32515
	IDC_UPARROW     = 32516
	IDC_SIZENWSE    = 32642
	IDC_SIZENESW    = 32643
	IDC_SIZEWE      = 32644
	IDC_SIZENS      = 32645
	IDC_SIZEALL     = 32646
	IDC_NO          = 32648
	IDC_HAND        = 32649
	IDC_APPSTARTING = 32650
	IDC_HELP        = 32651
	IDC_ICON        = 32641
	IDC_SIZE        = 32640
)

type CREATESTRUCT struct {
	CreateParams uintptr
	Instance     Handle
	Menu         Handle
	Parent       Handle
	Cy           int32
	Cx           int32
	Y            int32
	X            int32
	Style        int32
	Name         uintptr
	ClassName    uintptr
	ExStyle      uint32
}

type POINT struct {
	X int32
	Y int32
}

// RECT.
type RECT struct {
	Left   uint32
	Top    uint32
	Right  uint32
	Bottom uint32
}

// PAINTSTRUCT.
type PAINTSTRUCT struct {
	Hdc         Handle // HDC
	FErase      int32  // BOOL
	RcPaint     RECT   //
	FRestore    int32  // BOOL
	FIncUpdate  int32  // BOOL
	RegReserved [32]byte
}

// Draw Text constants.
const (
	DT_TOP                  = 0x00000000
	DT_LEFT                 = 0x00000000
	DT_CENTER               = 0x00000001
	DT_RIGHT                = 0x00000002
	DT_VCENTER              = 0x00000004
	DT_BOTTOM               = 0x00000008
	DT_WORDBREAK            = 0x00000010
	DT_SINGLELINE           = 0x00000020
	DT_EXPANDTABS           = 0x00000040
	DT_TABSTOP              = 0x00000080
	DT_NOCLIP               = 0x00000100
	DT_EXTERNALLEADING      = 0x00000200
	DT_CALCRECT             = 0x00000400
	DT_NOPREFIX             = 0x00000800
	DT_INTERNAL             = 0x00001000
	DT_EDITCONTROL          = 0x00002000
	DT_PATH_ELLIPSIS        = 0x00004000
	DT_END_ELLIPSIS         = 0x00008000
	DT_MODIFYSTRING         = 0x00010000
	DT_RTLREADING           = 0x00020000
	DT_WORD_ELLIPSIS        = 0x00040000
	DT_NOFULLWIDTHCHARBREAK = 0x00080000
	DT_HIDEPREFIX           = 0x00100000
	DT_PREFIXONLY           = 0x00200000
)

type DRAWTEXTPARAMS struct {
	Size          uint32
	ITabLength    int32
	ILeftMargin   int32
	IRightMargin  int32
	UiLengthDrawn uint32
}

// Bitmap compression.
const (
	BI_RGB       = 0
	BI_RLE8      = 1
	BI_RLE4      = 2
	BI_BITFIELDS = 3
	BI_JPEG      = 4
	BI_PNG       = 5
)

// BITMAPINFOHEADER.
type BITMAPINFOHEADER struct {
	Size          uint32 // DWORD
	Width         int32  // LONG
	Height        int32  // LONG
	Planes        uint16 // WORD
	BitCount      uint16 // WORD
	Compression   uint32 // DWORD
	SizeImage     uint32 // DWORD
	XPelsPerMeter int32  // LONG
	YPelsPerMeter int32  // LONG
	ClrUsed       uint32 // DWORD
	ClrImportant  uint32 // DWORD
}

// RGBQUAD.
type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

// BITMAPINFO.
type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors *RGBQUAD
}

// Bitmap color table usage.
const (
	DIB_RGB_COLORS = 0
	DIB_PAL_COLORS = 1
)

// Ternary raster operations
const (
	SRCCOPY        = 0x00CC0020
	SRCPAINT       = 0x00EE0086
	SRCAND         = 0x008800C6
	SRCINVERT      = 0x00660046
	SRCERASE       = 0x00440328
	NOTSRCCOPY     = 0x00330008
	NOTSRCERASE    = 0x001100A6
	MERGECOPY      = 0x00C000CA
	MERGEPAINT     = 0x00BB0226
	PATCOPY        = 0x00F00021
	PATPAINT       = 0x00FB0A09
	PATINVERT      = 0x005A0049
	DSTINVERT      = 0x00550009
	BLACKNESS      = 0x00000042
	WHITENESS      = 0x00FF0062
	NOMIRRORBITMAP = 0x80000000
	CAPTUREBLT     = 0x40000000
)

// GdiAlphaBlend Arg
type BLENDFUNCTION struct {
	BlendOp             byte
	BlendFlags          byte
	SourceConstantAlpha byte
	AlphaFormat         byte
}

const (
	AC_SRC_OVER             = 0x00
	AC_SRC_ALPHA            = 0x01
	AC_SRC_NO_PREMULT_ALPHA = 0x01
	AC_SRC_NO_ALPHA         = 0x02
	AC_DST_NO_PREMULT_ALPHA = 0x10
	AC_DST_NO_ALPHA         = 0x20
)

const (
	RDW_ERASE           = 4
	RDW_FRAME           = 1024
	RDW_INTERNALPAINT   = 2
	RDW_INVALIDATE      = 1
	RDW_NOERASE         = 32
	RDW_NOFRAME         = 2048
	RDW_NOINTERNALPAINT = 16
	RDW_VALIDATE        = 8
	RDW_ERASENOW        = 512
	RDW_UPDATENOW       = 256
	RDW_ALLCHILDREN     = 128
	RDW_NOCHILDREN      = 64
)
