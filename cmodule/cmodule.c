/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include "cmodule.h"

typedef struct {int64_t err1, err2; const char *err_str; uintptr_t *data, *ext; size_t length, index; int32_t pass;} cmodule_params_t;
typedef void (*cmodule_proc_t)(cmodule_params_t *params);

void vbsw_cmodule_alloc(uintptr_t **const data, const size_t length) {
	const size_t size = length*sizeof(void*);
	void *const data_new = malloc(size);
	memset(data_new, 0, size);
	*data = (uintptr_t*)data_new;
}

void vbsw_cmodule_proc(cmodule_proc_params_t *const proc_params) {
	cmodule_params_t params = {0, 0, 0, &proc_params->data[proc_params->length], &proc_params->data[proc_params->length*2], proc_params->length, 0, 0};
	// main
	while (params.pass < proc_params->passes) {
		// forward
		for (params.index = 0; params.index < proc_params->length && !params.err1; params.index++) {
			if (proc_params->data[params.index]) {
				proc_params->err_idx = params.index;
				((cmodule_proc_t)proc_params->data[params.index])(&params);
			}
		}
		// backwards
		if (!params.err1 && ++params.pass < proc_params->passes) {
			for (params.index = proc_params->length - 1; params.index >= 0 && !params.err1; params.index--) {
				if (proc_params->data[params.index]) {
					proc_params->err_idx = params.index;
					((cmodule_proc_t)proc_params->data[params.index])(&params);
				}
			}
		}
		if (!params.err1)
			params.pass++;
		else
			break;
	}
	// error handling
	if (params.err1) {
		params.pass = -(params.pass + 1);
		for (params.index = proc_params->length - 1; params.index >= 0; params.index--) {
			if (proc_params->data[params.index]) {
				((cmodule_proc_t)proc_params->data[params.index])(&params);
			}
		}
	}
}

void cmodule_free(uintptr_t *const data) {
	free((void*)data);
}
