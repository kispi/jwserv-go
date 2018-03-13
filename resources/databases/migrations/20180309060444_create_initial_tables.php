<?php

use Phinx\Migration\AbstractMigration;

class CreateInitialTables extends AbstractMigration
{
    protected $tables = [
        "congregations",
        "service_records",
        "users"
    ];

    // 'boolean', ('enum', ['values' => ['female', 'male']])
    public function change() {
        $this->table('congregations', ['id' => false, 'primary_key' => 'id'])
        ->addColumn('id', 'biginteger', ['identity' => true])
        ->addColumn('name', 'string', ['null' => true])
        ->addColumn('number', 'string', ['null' => true])
        ->save();

        $this->table('users', ['id' => false, 'primary_key' => 'id'])
        ->addColumn('id', 'biginteger', ['identity' => true])
        ->addColumn('nickname', 'string', ['null' => true])
        ->addColumn('congregation_id', 'biginteger', ['null' => true])
        ->addColumn('email', 'string', ['null' => true])
        ->addColumn('phone', 'string', ['null' => true])
        ->addColumn('name', 'string', ['null' => true])
        ->addColumn('password', 'string', ['null' => true])
        ->addColumn('auth', 'enum', ['values' => ['a', 'w', 'r'], 'default' => 'r'])
        ->addColumn('last_activity', 'timestamp', ['null' => true])
        ->addForeignKey('congregation_id', 'congregations', 'id', array('delete'=> 'SET NULL', 'update'=> 'CASCADE'))
        ->save();

        $this->table('service_records', ['id' => false, 'primary_key' => 'id'])
        ->addColumn('id', 'biginteger', ['identity' => true])
        ->addColumn('area', 'string', ['null' => true])
        ->addColumn('started_at', 'timestamp', ['null' => true])
        ->addColumn('ended_at', 'timestamp', ['null' => true])
        ->addColumn('congregation_id', 'biginteger', ['null' => true])
        ->addColumn('leader_name', 'string', ['null' => true])
        ->addColumn('recorder_id', 'biginteger', ['null' => true])
        ->addColumn('memo', 'text', ['null' => true])
        ->addForeignKey('congregation_id', 'congregations', 'id', array('delete'=> 'SET NULL', 'update'=> 'CASCADE'))
        ->addForeignKey('recorder_id', 'users', 'id', array('delete'=> 'SET NULL', 'update'=> 'CASCADE'))
        ->save();
    }
}