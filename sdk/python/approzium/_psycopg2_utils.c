#define PY_SSIZE_T_CLEAN
#include <Python.h>

struct connectionObject {
    PyObject_HEAD

    pthread_mutex_t lock;   /* the global connection lock */

    char *dsn;              /* data source name */
    char *error;            /* temporarily stored error before raising */
    char *encoding;         /* current backend encoding */

    long int closed;          /* 1 means connection has been closed;
                                                          2 that something horrible happened */
    long int mark;            /* number of commits/rollbacks done so far */
    int status;               /* status of the connection */
    void *tpc_xid;       /* Transaction ID in two-phase commit */
    long int async;           /* 1 means the connection is async */
};

static PyObject *psycopg2_utils_set_sync(PyObject *self, PyObject *args) {
    PyObject *object;
    
    if (!PyArg_ParseTuple(args, "O", &object))
        return NULL;

    struct connectionObject *connection = (struct connectionObject *)object;
    connection->async = 0;
    return PyLong_FromLong(0);
}

static PyMethodDef _Psycopg2UtilsMethods[] = {
    {"set_sync", psycopg2_utils_set_sync, METH_VARARGS,
        "Set Psycopg2 conection to sync"},
    {NULL, NULL, 0, NULL}
};

static struct PyModuleDef _psycopg2_utilsmodule = {
    PyModuleDef_HEAD_INIT,
    "_psycopg2_utils",
    NULL,
    -1,
    _Psycopg2UtilsMethods
};

PyMODINIT_FUNC PyInit__psycopg2_utils(void)
{
    return PyModule_Create(&_psycopg2_utilsmodule);
}
