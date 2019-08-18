#include <mysql.h>

#if defined(_WIN32) || defined(_WIN64) || defined(__WIN32__) || defined(WIN32)
    #define DLLEXP __declspec(dllexport)
#else
    #define DLLEXP
#endif

DLLEXP long long udf_add(UDF_INIT *initid, UDF_ARGS *args, char *is_null, char *error);
DLLEXP bool udf_add_init(UDF_INIT *initid, UDF_ARGS *args, char *message);
DLLEXP void udf_add_deinit(UDF_INIT *initid);

long long udf_add(UDF_INIT *initid, UDF_ARGS *args, char *is_null, char *error)
{
        int a = *((long long *)args->args[0]);
        int b = *((long long *)args->args[1]);

        return a + b;
}

bool udf_add_init(UDF_INIT *initid, UDF_ARGS *args, char *message)
{
        return 0;
}

void udf_add_deinit(UDF_INIT *initid)
{
        //
}