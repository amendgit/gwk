#ifndef EVENTS_H
#define EVENTS_H

#ifdef __cplusplus
extern "C" {
#endif

enum {
	WindowEvent_Resize = 511,
	WindowEvent_Move   = 512,

	WindowEvent_Close   = 521,
	WindowEvent_Destroy = 522,

	WindowEvent_Minimize = 531,
	WindowEvent_Maximize = 532,
	WindowEvent_Restore  = 533,

	WindowEvent_FocusMin            = 541,
	WindowEvent_FocusLost           = 541,
	WindowEvent_FocusGained         = 542,
	WindowEvent_FocusGainedForward  = 543,
	WindowEvent_FocusGainedBackward = 544,
	WindowEvent_FocusMax            = 544,

	WindowEvent_FocusDisabled = 545,
	WindowEvent_FocusUngrab   = 546,

	WindowEvent_InitAccessibility = 551,
};

const char* WindowEventName(int eventType);


enum {
	MouseEvent_ButtonNone  = 211,
	MouseEvent_ButtonLeft  = 212,
	MouseEvent_ButtonRight = 213,
	MouseEvent_ButtonOther = 214,

	MouseEvent_Down  = 221,
	MouseEvent_Up    = 222,
	MouseEvent_Drag  = 223,
	MouseEvent_Move  = 224,
	MouseEvent_Enter = 225,
	MouseEvent_Exit  = 226,
	MouseEvent_Click = 227, // synthetic

	// Artificial Whell event type.
	// This kind of mouse event is NEVER sent to app.
	// The app must listen to Scroll events instead.
	// This identifier is required for internal purposes.
	MouseEvent_Wheel = 228,
};

enum {
	// Key event type.
	KeyEvent_Press = 111,
	KeyEvent_Release = 112,
	KeyEvent_Typed = 113,

	// Key event modifier flags.
	// CONTROL/WINDOWS and OPTION/ALT are equal, because they
	// are mapped to each other on Mac/Windows.
	KeyEvent_ModifierNone     = 0,
	KeyEvent_ModifierShift    = 1 << 0,
	KeyEvent_ModifierFunction = 1 << 1,
	KeyEvent_ModifierControl  = 1 << 2,

	KeyEvent_ModifierOption = 1 << 3,
	KeyEvent_ModifierAlt    = 1 << 3,

	// The fellowing should be named Meta perhaps?
	KeyEvent_ModifierCommand = 1 << 4,
	KeyEvent_ModifierWindows = 1 << 4,
	KeyEvent_ModifierMeta    = 1 << 4,

	// Mouse Buttons
	KeyEvent_ModifierButtonPrimary   = 1 << 5,
	KeyEvent_ModifierButtonSecondary = 1 << 6,
	KeyEvent_ModifierButtonMiddle    = 1 << 7,

	// Key event key codes.
	KeyCode_Undefined = 0x0,

	// Misc
	KeyCode_Enter       = '\n',
	KeyCode_Backspace   = '\b',
	KeyCode_Tab         = '\t',
	KeyCode_Clear       = 0x0C,
	KeyCode_Pause       = 0x13,
	KeyCode_Escape      = 0x1B,
	KeyCode_Space       = 0x20,
	KeyCode_Delete      = 0x7F,
	KeyCode_PrintScreen = 0x9A,
	KeyCode_Insert      = 0x9B,
	KeyCode_Help        = 0x9C,

	// Modifiers
	KeyCode_Shift       = 0x10,
	KeyCode_Control     = 0x11,
	KeyCode_Alt         = 0x12,
	KeyCode_AltGraph    = 0xFF7E,
	KeyCode_Windows     = 0x020C,
	KeyCode_ContextMenu = 0x020D,
	KeyCode_CapsLock    = 0x14,
	KeyCode_NumLock     = 0x90,
	KeyCode_ScrollLock  = 0x91,
	KeyCode_Command     = 0x0300,

	// Navigation Keys
	KeyCode_PageUp   = 0x21,
	KeyCode_PageDown = 0x22,
	KeyCode_End      = 0x23,
	KeyCode_Home     = 0x24,
	KeyCode_Left     = 0x25,
	KeyCode_Up       = 0x26,
	KeyCode_Right    = 0x27,
	KeyCode_Down     = 0x28,

	// Misc 2
	KeyCode_Comma            = 0x2C,
	KeyCode_Minus            = 0x2D,
	KeyCode_Period           = 0x2E,
	KeyCode_Slash            = 0x2F,
	KeyCode_Semicolon        = 0x3B,
	KeyCode_Equals           = 0x3D,
	KeyCode_OpenBracket      = 0x5B,
	KeyCode_Backslash        = 0x5C,
	KeyCode_CloseBracket     = 0x5D,
	KeyCode_Multiply         = 0x6A,
	KeyCode_Add              = 0x6B,
	KeyCode_Separator        = 0x6C,
	KeyCode_Subtract         = 0x6D,
	KeyCode_Decimal          = 0x6E,
	KeyCode_Divide           = 0x6F,
	KeyCode_Ampersand        = 0x96,
	KeyCode_Asterisk         = 0x97,
	KeyCode_DoubleQuote      = 0x98,
	KeyCode_Less             = 0x99,
	KeyCode_Greater          = 0xA0,
	KeyCode_BraceLeft        = 0xA1,
	KeyCode_BraceRight       = 0xA2,
	KeyCode_BackQuote        = 0xC0,
	KeyCode_Quote            = 0xDE,
	KeyCode_At               = 0x0200,
	KeyCode_Colon            = 0x0201,
	KeyCode_Circumflex       = 0x0202,
	KeyCode_Dollar           = 0x0203,
	KeyCode_EuroSign         = 0x0204,
	KeyCode_Exclamation      = 0x0205,
	KeyCode_InvExclamation   = 0x0206,
	KeyCode_LeftParenthesis  = 0x0207,
	KeyCode_NumberSign       = 0x0208,
	KeyCode_Plus             = 0x0209,
	KeyCode_RightParenthesis = 0x020A,
	KeyCode_Underscore       = 0x020B,

	// Numberic keys.
	KeyCode_0 = 0x30,
	KeyCode_1 = 0x31,
	KeyCode_2 = 0x32,
	KeyCode_3 = 0x33,
	KeyCode_4 = 0x34,
	KeyCode_5 = 0x35,
	KeyCode_6 = 0x36,
	KeyCode_7 = 0x37,
	KeyCode_8 = 0x38,
	KeyCode_9 = 0x39,

	// Alpha keys.
	KeyCode_A = 0x41,
	KeyCode_B = 0x42,
	KeyCode_C = 0x43,
	KeyCode_D = 0x44,
	KeyCode_E = 0x45,
	KeyCode_F = 0x46,
	KeyCode_G = 0x47,
	KeyCode_H = 0x48,
	KeyCode_I = 0x49,
	KeyCode_J = 0x4A,
	KeyCode_K = 0x4B,
	KeyCode_L = 0x4C,
	KeyCode_M = 0x4D,
	KeyCode_N = 0x4E,
	KeyCode_O = 0x4F,
	KeyCode_P = 0x50,
	KeyCode_Q = 0x51,
	KeyCode_R = 0x52,
	KeyCode_S = 0x53,
	KeyCode_T = 0x54,
	KeyCode_U = 0x55,
	KeyCode_V = 0x56,
	KeyCode_W = 0x57,
	KeyCode_X = 0x58,
	KeyCode_Y = 0x59,
	KeyCode_Z = 0x60,

	// Numpad keys.
	KeyCode_Numpad0 = 0x60,
	KeyCode_Numpad1 = 0x61,
	KeyCode_Numpad2 = 0x62,
	KeyCode_Numpad3 = 0x63,
	KeyCode_Numpad4 = 0x64,
	KeyCode_Numpad5 = 0x65,
	KeyCode_Numpad6 = 0x66,
	KeyCode_Numpad7 = 0x67,
	KeyCode_Numpad8 = 0x68,
	KeyCode_Numpad9 = 0x69,

	// Function keys.
	KeyCode_F1  = 0x70,
	KeyCode_F2  = 0x71,
	KeyCode_F3  = 0x72,
	KeyCode_F4  = 0x73,
	KeyCode_F5  = 0x74,
	KeyCode_F6  = 0x75,
	KeyCode_F7  = 0x76,
	KeyCode_F8  = 0x77,
	KeyCode_F9  = 0x78,
	KeyCode_F10 = 0x79,
	KeyCode_F11 = 0x7A,
	KeyCode_F12 = 0x7B,
	KeyCode_F13 = 0xF000,
	KeyCode_F14 = 0xF001,
	KeyCode_F15 = 0xF002,
	KeyCode_F16 = 0xF003,
	KeyCode_F17 = 0xF004,
	KeyCode_F18 = 0xF005,
	KeyCode_F19 = 0xF006,
	KeyCode_F20 = 0xF007,
	KeyCode_F21 = 0xF008,
	KeyCode_F22 = 0xF009,
	KeyCode_F23 = 0xF00A,
	KeyCode_F24 = 0xF00B,
};

int KeyCodeFromChar(char c);
const char *KeyEventName(int eventType);

#ifdef __cplusplus
}
#endif

#endif
