#ifndef CMODULE_H
#define CMODULE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
	int64_t err1, err2;
	char *err_str;
	uintptr_t *data;
	size_t err_idx, length;
	int32_t passes;
} cmodule_proc_params_t;

extern void cmodule_alloc(uintptr_t **data, size_t length);
extern void cmodule_proc(cmodule_proc_params_t *params);
extern void cmodule_free(uintptr_t *data);

#ifdef __cplusplus
}
#endif

#endif /* CMODULE_H */
