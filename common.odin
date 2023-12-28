package invaders

import "core:encoding/json"
import "core:os"
import "core:strings"

Window :: struct {
	width:      i32,
	height:     i32,
	fullscreen: bool,
}

CONFIG_NAME :: "config.json"

InitWindowError :: union {
	json.Unmarshal_Error,
	json.Marshal_Error,
}

init_window :: proc(allocator := context.allocator) -> (window: Window, err: InitWindowError) {
	file_data, err_file := os.read_entire_file_from_filename(CONFIG_NAME, allocator)
	if !err_file {
		default_data := json.marshal(
			Window{width = 800, height = 450, fullscreen = false}, 
            json.Marshal_Options {
                pretty = true
            }, 
            allocator
		) or_return

		os.write_entire_file(CONFIG_NAME, default_data)
	}

	defer delete(file_data)

    json.unmarshal(file_data, &window, json.DEFAULT_SPECIFICATION, allocator) or_return

	return window, nil
}
