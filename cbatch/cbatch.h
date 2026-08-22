#ifndef CMODULE_H
#define CMODULE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
	int64_t err1, err2;
	const char *err_str;
	void **data;
	size_t err_idx, length;
	int32_t passes;
} cbatch_proc_params_t;

extern void cbatch_alloc(void ***data, size_t total_length);
extern void cbatch_proc(cbatch_proc_params_t *params);
extern void cbatch_free(void **data);

#ifdef __cplusplus
}
#endif

#endif /* CMODULE_H */
