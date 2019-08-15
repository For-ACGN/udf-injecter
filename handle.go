package injecter

type Handle interface {
	// Query
	// mysql> show variables like '%version_%';
	// +-------------------------+------------------------------+
	// | Variable_name           | Value                        |
	// +-------------------------+------------------------------+
	// | slave_type_conversions  |                              |
	// | version_comment         | MySQL Community Server - GPL |
	// | version_compile_machine | x86_64                       |
	// | version_compile_os      | Win64                        |
	// | version_compile_zlib    | 1.2.11                       |
	// +-------------------------+------------------------------+
	// result, _ := Query("show variables like '%version_%'")
	// result[2]["Variable_name"] = "x86_64"
	// Query result length must > 0
	Query(query string, args ...interface{}) ([]map[string]string, error)
	Exec(query string, args ...interface{}) error
}
