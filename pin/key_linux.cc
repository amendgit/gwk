#include "key_linux.h"
#include "gtkcompat_linux.h"
#include "events.h"

static gboolean g_isKeyMapInitialized = FALSE;
static GHashTable *g_keyMap = NULL;

static void KeyMapInsert(gint key, gint val) {
    g_hash_table_insert(g_keyMap, GINT_TO_POINTER(key), GINT_TO_POINTER(val));
}

static void KeyMapInit() {
    if (g_isKeyMapInitialized) return ;
    g_isKeyMapInitialized = TRUE;

    g_keyMap = g_hash_table_new(g_direct_hash, g_direct_equal);
    KeyMapInsert(GDK_KEY_CONSTANT(Return),    kKeyCodeEnter);
    KeyMapInsert(GDK_KEY_CONSTANT(BackSpace), kKeyCodeBackspace);
    KeyMapInsert(GDK_KEY_CONSTANT(Tab),       kKeyCodeTab);
    KeyMapInsert(GDK_KEY_CONSTANT(Clear),     kKeyCodeClear);
    KeyMapInsert(GDK_KEY_CONSTANT(Pause),     kKeyCodePause);
    KeyMapInsert(GDK_KEY_CONSTANT(Escape),    kKeyCodeEscape);
    KeyMapInsert(GDK_KEY_CONSTANT(space),     kKeyCodeSpace);
    KeyMapInsert(GDK_KEY_CONSTANT(Delete),    kKeyCodeDelete);
    KeyMapInsert(GDK_KEY_CONSTANT(Print),     kKeyCodePrintScreen);
    KeyMapInsert(GDK_KEY_CONSTANT(Insert),    kKeyCodeInsert);
    KeyMapInsert(GDK_KEY_CONSTANT(Help),      kKeyCodeHelp);

    KeyMapInsert(GDK_KEY_CONSTANT(Shift_L),     kKeyCodeShift);
    KeyMapInsert(GDK_KEY_CONSTANT(Shift_R),     kKeyCodeShift);
    KeyMapInsert(GDK_KEY_CONSTANT(Control_L),   kKeyCodeControl);
    KeyMapInsert(GDK_KEY_CONSTANT(Control_R),   kKeyCodeControl);
    KeyMapInsert(GDK_KEY_CONSTANT(Alt_L),       kKeyCodeAlt);
    KeyMapInsert(GDK_KEY_CONSTANT(Alt_R),       kKeyCodeAltGraph);
    KeyMapInsert(GDK_KEY_CONSTANT(Super_L),     kKeyCodeWindows);
    KeyMapInsert(GDK_KEY_CONSTANT(Super_R),      kKeyCodeWindows);
    KeyMapInsert(GDK_KEY_CONSTANT(Menu),        kKeyCodeContextMenu);
    KeyMapInsert(GDK_KEY_CONSTANT(Caps_Lock),   kKeyCodeCapsLock);
    KeyMapInsert(GDK_KEY_CONSTANT(Num_Lock),    kKeyCodeNumLock);
    KeyMapInsert(GDK_KEY_CONSTANT(Scroll_Lock), kKeyCodeScrollLock);

    KeyMapInsert(GDK_KEY_CONSTANT(Page_Up),   kKeyCodePageUp);
    KeyMapInsert(GDK_KEY_CONSTANT(Prior),     kKeyCodePageUp);
    KeyMapInsert(GDK_KEY_CONSTANT(Page_Down), kKeyCodePageDown);
    KeyMapInsert(GDK_KEY_CONSTANT(Next),      kKeyCodePageDown);
    KeyMapInsert(GDK_KEY_CONSTANT(End),       kKeyCodeEnd);
    KeyMapInsert(GDK_KEY_CONSTANT(Home),      kKeyCodeHome);
    KeyMapInsert(GDK_KEY_CONSTANT(Left),      kKeyCodeLeft);
    KeyMapInsert(GDK_KEY_CONSTANT(Right),     kKeyCodeRight);
    KeyMapInsert(GDK_KEY_CONSTANT(Up),        kKeyCodeUp);
    KeyMapInsert(GDK_KEY_CONSTANT(Down),      kKeyCodeDown);

    KeyMapInsert(GDK_KEY_CONSTANT(comma),        kKeyCodeComma);
    KeyMapInsert(GDK_KEY_CONSTANT(minus),        kKeyCodeMinus);
    KeyMapInsert(GDK_KEY_CONSTANT(period),       kKeyCodePeriod);
    KeyMapInsert(GDK_KEY_CONSTANT(slash),        kKeyCodeSlash);
    KeyMapInsert(GDK_KEY_CONSTANT(semicolon),    kKeyCodeSemicolon);
    KeyMapInsert(GDK_KEY_CONSTANT(equal),        kKeyCodeEquals);
    KeyMapInsert(GDK_KEY_CONSTANT(bracketleft),  kKeyCodeLeftBracket);
    KeyMapInsert(GDK_KEY_CONSTANT(bracketright), kKeyCodeRightBracket);
    KeyMapInsert(GDK_KEY_CONSTANT(backslash),    kKeyCodeBackslash);
    KeyMapInsert(GDK_KEY_CONSTANT(bar),          kKeyCodeBackslash);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Multiply),  kKeyCodeMultiply);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Add),       kKeyCodeAdd);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Separator), kKeyCodeSeparator);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Subtract),  kKeyCodeSubtract);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Decimal),   kKeyCodeDecimal);

    KeyMapInsert(GDK_KEY_CONSTANT(apostrophe), kKeyCodeQuote);
    KeyMapInsert(GDK_KEY_CONSTANT(grave),      kKeyCodeBackQuote);

    KeyMapInsert(GDK_KEY_CONSTANT(ampersand),   kKeyCodeAmpersand);
    KeyMapInsert(GDK_KEY_CONSTANT(asterisk),    kKeyCodeAsterisk);
    KeyMapInsert(GDK_KEY_CONSTANT(quotedbl),    kKeyCodeDoubleQuote);
    KeyMapInsert(GDK_KEY_CONSTANT(less),        kKeyCodeLess);
    KeyMapInsert(GDK_KEY_CONSTANT(greater),     kKeyCodeGreater);
    KeyMapInsert(GDK_KEY_CONSTANT(braceleft),   kKeyCodeLeftBrace);
    KeyMapInsert(GDK_KEY_CONSTANT(braceright),  kKeyCodeRightBrace);
    KeyMapInsert(GDK_KEY_CONSTANT(at),          kKeyCodeAt);
    KeyMapInsert(GDK_KEY_CONSTANT(colon),       kKeyCodeColon);
    KeyMapInsert(GDK_KEY_CONSTANT(asciicircum), kKeyCodeCircumflex);
    KeyMapInsert(GDK_KEY_CONSTANT(dollar),      kKeyCodeDollar);
    KeyMapInsert(GDK_KEY_CONSTANT(EuroSign),    kKeyCodeEuroSign);
    KeyMapInsert(GDK_KEY_CONSTANT(exclam),      kKeyCodeExclamation);
    KeyMapInsert(GDK_KEY_CONSTANT(exclamdown),  kKeyCodeInvExclamation);
    KeyMapInsert(GDK_KEY_CONSTANT(parenleft),   kKeyCodeLeftParenthesis);
    KeyMapInsert(GDK_KEY_CONSTANT(parenright),  kKeyCodeRightParenthesis);
    KeyMapInsert(GDK_KEY_CONSTANT(numbersign),  kKeyCodeNumberSign);
    KeyMapInsert(GDK_KEY_CONSTANT(plus),        kKeyCodePlus);
    KeyMapInsert(GDK_KEY_CONSTANT(underscore),  kKeyCodeUnderscore);

    KeyMapInsert(GDK_KEY_CONSTANT(0), kKeyCode0);
    KeyMapInsert(GDK_KEY_CONSTANT(1), kKeyCode1);
    KeyMapInsert(GDK_KEY_CONSTANT(2), kKeyCode2);
    KeyMapInsert(GDK_KEY_CONSTANT(3), kKeyCode3);
    KeyMapInsert(GDK_KEY_CONSTANT(4), kKeyCode4);
    KeyMapInsert(GDK_KEY_CONSTANT(5), kKeyCode5);
    KeyMapInsert(GDK_KEY_CONSTANT(6), kKeyCode6);
    KeyMapInsert(GDK_KEY_CONSTANT(7), kKeyCode7);
    KeyMapInsert(GDK_KEY_CONSTANT(8), kKeyCode8);
    KeyMapInsert(GDK_KEY_CONSTANT(9), kKeyCode9);

    KeyMapInsert(GDK_KEY_CONSTANT(a), kKeyCodeA);
    KeyMapInsert(GDK_KEY_CONSTANT(b), kKeyCodeB);
    KeyMapInsert(GDK_KEY_CONSTANT(c), kKeyCodeC);
    KeyMapInsert(GDK_KEY_CONSTANT(d), kKeyCodeD);
    KeyMapInsert(GDK_KEY_CONSTANT(e), kKeyCodeE);
    KeyMapInsert(GDK_KEY_CONSTANT(f), kKeyCodeF);
    KeyMapInsert(GDK_KEY_CONSTANT(g), kKeyCodeG);
    KeyMapInsert(GDK_KEY_CONSTANT(h), kKeyCodeH);
    KeyMapInsert(GDK_KEY_CONSTANT(i), kKeyCodeI);
    KeyMapInsert(GDK_KEY_CONSTANT(j), kKeyCodeJ);
    KeyMapInsert(GDK_KEY_CONSTANT(k), kKeyCodeK);
    KeyMapInsert(GDK_KEY_CONSTANT(l), kKeyCodeL);
    KeyMapInsert(GDK_KEY_CONSTANT(m), kKeyCodeM);
    KeyMapInsert(GDK_KEY_CONSTANT(n), kKeyCodeN);
    KeyMapInsert(GDK_KEY_CONSTANT(o), kKeyCodeO);
    KeyMapInsert(GDK_KEY_CONSTANT(p), kKeyCodeP);
    KeyMapInsert(GDK_KEY_CONSTANT(q), kKeyCodeQ);
    KeyMapInsert(GDK_KEY_CONSTANT(r), kKeyCodeR);
    KeyMapInsert(GDK_KEY_CONSTANT(s), kKeyCodeS);
    KeyMapInsert(GDK_KEY_CONSTANT(t), kKeyCodeT);
    KeyMapInsert(GDK_KEY_CONSTANT(u), kKeyCodeU);
    KeyMapInsert(GDK_KEY_CONSTANT(v), kKeyCodeV);
    KeyMapInsert(GDK_KEY_CONSTANT(w), kKeyCodeW);
    KeyMapInsert(GDK_KEY_CONSTANT(x), kKeyCodeX);
    KeyMapInsert(GDK_KEY_CONSTANT(y), kKeyCodeY);
    KeyMapInsert(GDK_KEY_CONSTANT(z), kKeyCodeZ);

    KeyMapInsert(GDK_KEY_CONSTANT(KP_0), kKeyCodeNumpad0);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_1), kKeyCodeNumpad1);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_2), kKeyCodeNumpad2);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_3), kKeyCodeNumpad3);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_4), kKeyCodeNumpad4);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_5), kKeyCodeNumpad5);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_6), kKeyCodeNumpad6);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_7), kKeyCodeNumpad7);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_8), kKeyCodeNumpad8);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_9), kKeyCodeNumpad9);

    KeyMapInsert(GDK_KEY_CONSTANT(KP_Enter),     kKeyCodeEnter);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Home),      kKeyCodeHome);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Left),      kKeyCodeLeft);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Up),        kKeyCodeUp);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Right),     kKeyCodeRight);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Down),      kKeyCodeDown);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Prior),     kKeyCodePageUp);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Page_Up),   kKeyCodePageUp);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Next),      kKeyCodePageDown);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Page_Down), kKeyCodePageDown);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_End),       kKeyCodeEnd);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Insert),    kKeyCodeInsert);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Delete),    kKeyCodeDelete);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Divide),    kKeyCodeDivide);
    KeyMapInsert(GDK_KEY_CONSTANT(KP_Begin),     kKeyCodeClear);

    KeyMapInsert(GDK_KEY_CONSTANT(F1),  kKeyCodeF1);
    KeyMapInsert(GDK_KEY_CONSTANT(F2),  kKeyCodeF2);
    KeyMapInsert(GDK_KEY_CONSTANT(F3),  kKeyCodeF3);
    KeyMapInsert(GDK_KEY_CONSTANT(F4),  kKeyCodeF4);
    KeyMapInsert(GDK_KEY_CONSTANT(F5),  kKeyCodeF5);
    KeyMapInsert(GDK_KEY_CONSTANT(F6),  kKeyCodeF6);
    KeyMapInsert(GDK_KEY_CONSTANT(F7),  kKeyCodeF7);
    KeyMapInsert(GDK_KEY_CONSTANT(F8),  kKeyCodeF8);
    KeyMapInsert(GDK_KEY_CONSTANT(F9),  kKeyCodeF9);
    KeyMapInsert(GDK_KEY_CONSTANT(F10), kKeyCodeF10);
    KeyMapInsert(GDK_KEY_CONSTANT(F11), kKeyCodeF11);
    KeyMapInsert(GDK_KEY_CONSTANT(F12), kKeyCodeF12);
}

int GdkKeyvalToGwk(guint keyval) {
    if (!g_isKeyMapInitialized) KeyMapInit();
    return GPOINTER_TO_INT(g_hash_table_lookup(g_keyMap, GINT_TO_POINTER(keyval)));
}

int GetGwkKey(GdkEventKey *e) {
    if (!g_isKeyMapInitialized) KeyMapInit();

    guint keyval;
    guint state = e->state & GDK_MOD2_MASK; // Numlock test
    gdk_keymap_translate_keyboard_state(gdk_keymap_get_default(),
        e->hardware_keycode, static_cast<GdkModifierType>(state), e->group,
        &keyval, NULL, NULL, NULL);

    int key = GPOINTER_TO_INT(g_hash_table_lookup(g_keyMap,
        GINT_TO_POINTER(keyval)));

    if (!key) {
        // We failed to find a keyval in our keymap, this may happen with
        // non-latin layouts(e.g. Cyrillic). So here we try to find a keyval
        // from a default layout (we assume that it is a US-like one).
         GdkKeymapKey kk;
         kk.keycode = e->hardware_keycode;
         kk.group = kk.level = 0;

         keyval = gdk_keymap_lookup_key(gdk_keymap_get_default(), &kk);

         key = GPOINTER_TO_INT(g_hash_table_lookup(g_keyMap,
             GINT_TO_POINTER(keyval)));
    }

    return key;
}

int GwkKeyToModifier(int gwkKey) {
    switch (gwkKey) {
        case kKeyCodeShift:
            return kKeyEventModifierShift;
        case kKeyCodeAlt:
        case kKeyCodeAltGraph:
            return kKeyEventModifierAlt;
        case kKeyCodeControl:
            return kKeyEventModifierControl;
        case kKeyCodeWindows:
            return kKeyEventModifierWindows;
        default:
            return 0;
    }
}

int GdkModifierMaskToGwk(guint mask) {
	int gwkMask = 0;
	if (mask & GDK_SHIFT_MASK) {
		gwkMask |= kKeyEventModifierShift;
	}
	if (mask & GDK_CONTROL_MASK) {
		gwkMask |= kKeyEventModifierControl;
	}
	if (mask & GDK_MOD1_MASK) {
		gwkMask |= kKeyEventModifierAlt;
	}
	if (mask & GDK_META_MASK) {
		gwkMask |= kKeyEventModifierAlt;
	}
	if (mask & GDK_BUTTON1_MASK) {
		gwkMask |= kKeyEventModifierButtonPrimary;
	}
	if (mask & GDK_BUTTON2_MASK) {
		gwkMask |= kKeyEventModifierButtonMiddle;
	}
	if (mask & GDK_BUTTON3_MASK) {
		gwkMask |= kKeyEventModifierButtonSecondary;
	}
	if (mask & GDK_SUPER_MASK) {
		gwkMask |= kKeyEventModifierWindows;
	}
	return gwkMask;
}

gint FindGdkKeyvalForGwkKeycode(int code) {
    gint result = -1;
    GHashTableIter iter;
    gpointer key, val;
    if (!g_isKeyMapInitialized) KeyMapInit();
    g_hash_table_iter_init(&iter, g_keyMap);
    while (g_hash_table_iter_next(&iter, &key, &val)) {
        if (code == GPOINTER_TO_INT(val)) {
            result == GPOINTER_TO_INT(key);
            break;
        }
    }
    return result;
}
