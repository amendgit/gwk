#ifndef EVENTS_H
#define EVENTS_H

#ifdef __cplusplus
extern "C" {
#endif

enum {
	kWindowEventResize = 511,
	kWindowEventMove   = 512,

	kWindowEventClose   = 521,
	kWindowEventDestroy = 522,

	kWindowEventMinimize = 531,
	kWindowEventMaximize = 532,
	kWindowEventRestore  = 533,

	kWindowEventFocusMin            = 541,
	kWindowEventFocusLost           = 541,
	kWindowEventFocusGained         = 542,
	kWindowEventFocusGainedForward  = 543,
	kWindowEventFocusGainedBackward = 544,
	kWindowEventFocusMax            = 544,

	kWindowEventFocusDisabled = 545,
	kWindowEventFocusUngrab   = 546,

	kWindowEventInitAccessibility = 551,
};

const char* WindowEventName(int eventType);


enum {
	kMouseEventButtonNone  = 211,
	kMouseEventButtonLeft  = 212,
	kMouseEventButtonRight = 213,
	kMouseEventButtonOther = 214,

	kMouseEventDown  = 221,
	kMouseEventUp    = 222,
	kMouseEventDrag  = 223,
	kMouseEventMove  = 224,
	kMouseEventEnter = 225,
	kMouseEventExit  = 226,
	kMouseEventClick = 227, // synthetic

	// Artificial Whell event type.
	// This kind of mouse event is NEVER sent to app.
	// The app must listen to Scroll events instead.
	// This identifier is required for internal purposes.
	kMouseEventWheel = 228,
};

enum {
	// Key event type.
	kKeyEventPress = 111,
	kKeyEventRelease = 112,
	kKeyEventTyped = 113,

	// Key event modifier flags.
	// CONTROL/WINDOWS and OPTION/ALT are equal, because they
	// are mapped to each other on Mac/Windows.
	kKeyEventModifierNone     = 0,
	kKeyEventModifierShift    = 1 << 0,
	kKeyEventModifierFunction = 1 << 1,
	kKeyEventModifierControl  = 1 << 2,

	kKeyEventModifierOption = 1 << 3,
	kKeyEventModifierAlt    = 1 << 3,

	// The fellowing should be named Meta perhaps?
	kKeyEventModifierCommand = 1 << 4,
	kKeyEventModifierWindows = 1 << 4,
	kKeyEventModifierMeta    = 1 << 4,

	// Mouse Buttons
	kKeyEventModifierButtonPrimary   = 1 << 5,
	kKeyEventModifierButtonSecondary = 1 << 6,
	kKeyEventModifierButtonMiddle    = 1 << 7,

	// Key event key codes.
	kKeyCodeUndefined = 0x0,

	// Misc
	kKeyCodeEnter       = '\n',
	kKeyCodeBackspace   = '\b',
	kKeyCodeTab         = '\t',
	kKeyCodeClear       = 0x0C,
	kKeyCodePause       = 0x13,
	kKeyCodeEscape      = 0x1B,
	kKeyCodeSpace       = 0x20,
	kKeyCodeDelete      = 0x7F,
	kKeyCodePrintScreen = 0x9A,
	kKeyCodeInsert      = 0x9B,
	kKeyCodeHelp        = 0x9C,

	// Modifiers
	kKeyCodeShift       = 0x10,
	kKeyCodeControl     = 0x11,
	kKeyCodeAlt         = 0x12,
	kKeyCodeAltGraph    = 0xFF7E,
	kKeyCodeWindows     = 0x020C,
	kKeyCodeContextMenu = 0x020D,
	kKeyCodeCapsLock    = 0x14,
	kKeyCodeNumLock     = 0x90,
	kKeyCodeScrollLock  = 0x91,
	kKeyCodeCommand     = 0x0300,

	// Navigation Keys
	kKeyCodePageUp   = 0x21,
	kKeyCodePageDown = 0x22,
	kKeyCodeEnd      = 0x23,
	kKeyCodeHome     = 0x24,
	kKeyCodeLeft     = 0x25,
	kKeyCodeUp       = 0x26,
	kKeyCodeRight    = 0x27,
	kKeyCodeDown     = 0x28,

	// Misc 2
	kKeyCodeComma            = 0x2C,
	kKeyCodeMinus            = 0x2D,
	kKeyCodePeriod           = 0x2E,
	kKeyCodeSlash            = 0x2F,
	kKeyCodeSemicolon        = 0x3B,
	kKeyCodeEquals           = 0x3D,
	kKeyCodeLeftBracket      = 0x5B,
	kKeyCodeBackslash        = 0x5C,
	kKeyCodeRightBracket     = 0x5D,
	kKeyCodeMultiply         = 0x6A,
	kKeyCodeAdd              = 0x6B,
	kKeyCodeSeparator        = 0x6C,
	kKeyCodeSubtract         = 0x6D,
	kKeyCodeDecimal          = 0x6E,
	kKeyCodeDivide           = 0x6F,
	kKeyCodeAmpersand        = 0x96,
	kKeyCodeAsterisk         = 0x97,
	kKeyCodeDoubleQuote      = 0x98,
	kKeyCodeLess             = 0x99,
	kKeyCodeGreater          = 0xA0,
	kKeyCodeLeftBrace        = 0xA1,
	kKeyCodeRightBrace       = 0xA2,
	kKeyCodeBackQuote        = 0xC0,
	kKeyCodeQuote            = 0xDE,
	kKeyCodeAt               = 0x0200,
	kKeyCodeColon            = 0x0201,
	kKeyCodeCircumflex       = 0x0202,
	kKeyCodeDollar           = 0x0203,
	kKeyCodeEuroSign         = 0x0204,
	kKeyCodeExclamation      = 0x0205,
	kKeyCodeInvExclamation   = 0x0206,
	kKeyCodeLeftParenthesis  = 0x0207,
	kKeyCodeNumberSign       = 0x0208,
	kKeyCodePlus             = 0x0209,
	kKeyCodeRightParenthesis = 0x020A,
	kKeyCodeUnderscore       = 0x020B,

	// Numberic keys.
	kKeyCode0 = 0x30,
	kKeyCode1 = 0x31,
	kKeyCode2 = 0x32,
	kKeyCode3 = 0x33,
	kKeyCode4 = 0x34,
	kKeyCode5 = 0x35,
	kKeyCode6 = 0x36,
	kKeyCode7 = 0x37,
	kKeyCode8 = 0x38,
	kKeyCode9 = 0x39,

	// Alpha keys.
	kKeyCodeA = 0x41,
	kKeyCodeB = 0x42,
	kKeyCodeC = 0x43,
	kKeyCodeD = 0x44,
	kKeyCodeE = 0x45,
	kKeyCodeF = 0x46,
	kKeyCodeG = 0x47,
	kKeyCodeH = 0x48,
	kKeyCodeI = 0x49,
	kKeyCodeJ = 0x4A,
	kKeyCodeK = 0x4B,
	kKeyCodeL = 0x4C,
	kKeyCodeM = 0x4D,
	kKeyCodeN = 0x4E,
	kKeyCodeO = 0x4F,
	kKeyCodeP = 0x50,
	kKeyCodeQ = 0x51,
	kKeyCodeR = 0x52,
	kKeyCodeS = 0x53,
	kKeyCodeT = 0x54,
	kKeyCodeU = 0x55,
	kKeyCodeV = 0x56,
	kKeyCodeW = 0x57,
	kKeyCodeX = 0x58,
	kKeyCodeY = 0x59,
	kKeyCodeZ = 0x60,

	// Numpad keys.
	kKeyCodeNumpad0 = 0x60,
	kKeyCodeNumpad1 = 0x61,
	kKeyCodeNumpad2 = 0x62,
	kKeyCodeNumpad3 = 0x63,
	kKeyCodeNumpad4 = 0x64,
	kKeyCodeNumpad5 = 0x65,
	kKeyCodeNumpad6 = 0x66,
	kKeyCodeNumpad7 = 0x67,
	kKeyCodeNumpad8 = 0x68,
	kKeyCodeNumpad9 = 0x69,

	// Function keys.
	kKeyCodeF1  = 0x70,
	kKeyCodeF2  = 0x71,
	kKeyCodeF3  = 0x72,
	kKeyCodeF4  = 0x73,
	kKeyCodeF5  = 0x74,
	kKeyCodeF6  = 0x75,
	kKeyCodeF7  = 0x76,
	kKeyCodeF8  = 0x77,
	kKeyCodeF9  = 0x78,
	kKeyCodeF10 = 0x79,
	kKeyCodeF11 = 0x7A,
	kKeyCodeF12 = 0x7B,
	kKeyCodeF13 = 0xF000,
	kKeyCodeF14 = 0xF001,
	kKeyCodeF15 = 0xF002,
	kKeyCodeF16 = 0xF003,
	kKeyCodeF17 = 0xF004,
	kKeyCodeF18 = 0xF005,
	kKeyCodeF19 = 0xF006,
	kKeyCodeF20 = 0xF007,
	kKeyCodeF21 = 0xF008,
	kKeyCodeF22 = 0xF009,
	kKeyCodeF23 = 0xF00A,
	kKeyCodeF24 = 0xF00B,
};

int KeyCodeFromChar(char c);
const char *KeyEventName(int eventType);

#ifdef __cplusplus
}
#endif

#endif
