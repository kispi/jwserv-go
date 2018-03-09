<?php
require __DIR__ . '/vendor/autoload.php';

/*
|--------------------------------------------------------------------------
| Load Environtment Files
|--------------------------------------------------------------------------
|
| Before start, we need to load the configuration from env files
|
*/

$env = new Dotenv\Dotenv(__DIR__);
$env->load();

/*
|--------------------------------------------------------------------------
| Return Configurations
|--------------------------------------------------------------------------
|
| Next, we need to return array configuration for phinx to execute
|
*/

return array(
    "paths"        => array(
        "migrations" => "./resources/databases/migrations",
        "seeds"      => "./resources/databases/seeds"
    ),
    "environments" => array(
        "default_migration_table" => "migrations",
        "default_database"        => "dev",
        "dev"                     => array(
            "adapter" => "mysql",
            "host"    => getenv('DB_HOST'),
            "name"    => getenv('DB_NAME'),
            "user"    => getenv('DB_USER'),
            "pass"    => getenv('DB_PASS'),
            "port"    => getenv('DB_PORT'),
	    "charset" => "utf8",
        )
    )
);
