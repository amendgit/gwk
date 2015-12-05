#ifndef GENERAL_LINUX_H
#define GENERAL_LINUX_H

#ifdef VERBOSE
#define LOG0(msg) {printf(msg);fflush(stdout);}
#define LOG1(msg, param) {printf(msg, param);fflush(stdout);}
#define LOG2(msg, param1, param2) {printf(msg, param1, param2);fflush(stdout);}
#define LOG3(msg, param1, param2, param3) {printf(msg, param1, param2, param3);fflush(stdout);}
#define LOG4(msg, param1, param2, param3, param4) {printf(msg, param1, param2, param3, param4);fflush(stdout);}
#define LOG5(msg, param1, param2, param3, param4, param5) {printf(msg, param1, param2, param3, param4, param5);fflush(stdout);}

#define LOG_STRING_ARRAY(env, array) dump_jstring_array(env, array);

#define ERROR0(msg) {fprintf(stderr, msg);fflush(stderr);}
#define ERROR1(msg, param) {fprintf(stderr, msg, param);fflush(stderr);}
#define ERROR2(msg, param1, param2) {fprintf(stderr, msg, param1, param2);fflush(stderr);}
#define ERROR3(msg, param1, param2, param3) {fprintf(stderr, msg, param1, param2, param3);fflush(stderr);}
#define ERROR4(msg, param1, param2, param3, param4) {fprintf(stderr, msg, param1, param2, param3, param4);fflush(stderr);}
#else
#define LOG0(msg)
#define LOG1(msg, param)
#define LOG2(msg, param1, param2)
#define LOG3(msg, param1, param2, param3)
#define LOG4(msg, param1, param2, param3, param4)
#define LOG5(msg, param1, param2, param3, param4, param5)

#define LOG_STRING_ARRAY(env, array)

#define ERROR0(msg)
#define ERROR1(msg, param)
#define ERROR2(msg, param1, param2)
#define ERROR3(msg, param1, param2, param3)
#define ERROR4(msg, param1, param2, param3, param4)
#endif

#endif /* VERBOSE */
