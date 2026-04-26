#ifndef VBSW_CMODULE_H
#define VBSW_CMODULE_H

#ifdef __cplusplus
extern "C" {
#endif

extern void vbsw_cmodule_alloc_buffer(void ***data, size_t *data_len, void **data_old, size_t modules_len);
extern void vbsw_cmodule_proc(void **data, size_t modules_len, int passes, size_t *err_idx, long long *err1, long long *err2, char **err_str);
extern void vbsw_cmodule_rm(void **data, size_t modules_len, int passes, size_t *err_idx, long long *err1, long long *err2, char **err_str);

#ifdef __cplusplus
}
#endif

#endif /* VBSW_CMODULE_H */
