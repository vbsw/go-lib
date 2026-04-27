#ifndef VBSW_CMODULE_H
#define VBSW_CMODULE_H

#ifdef __cplusplus
extern "C" {
#endif

extern void vbsw_cmodule_alloc_buffer(void ***data, int32_t *data_len, int32_t *data_size, int32_t mod_len_new);
extern void vbsw_cmodule_proc(void **data, int32_t data_len, int32_t *data_size, int passes, int32_t *err_idx, int64_t *err1, int64_t *err2, char **err_str);
extern void vbsw_cmodule_free(void **data, int32_t data_len);

#ifdef __cplusplus
}
#endif

#endif /* VBSW_CMODULE_H */
