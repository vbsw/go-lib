#ifndef CMODULE_H
#define CMODULE_H

#ifdef __cplusplus
extern "C" {
#endif

void aca_batch_alloc(void ***data, int_fast32_t total_len);
void aca_batch_run(void **data, int_fast32_t *err1, int_fast32_t *err2, int_fast32_t *err_idx, char **err_str, int_fast32_t len, int_fast32_t passes);
void aca_batch_free(void **data);

#ifdef __cplusplus
}
#endif

#endif /* CMODULE_H */
