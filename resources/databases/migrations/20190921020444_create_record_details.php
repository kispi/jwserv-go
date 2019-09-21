<?php

use Phinx\Migration\AbstractMigration;

class CreateRecordDetails extends AbstractMigration
{
    public function change() {
        $table = 'record_details';

        $this->table($table, ['id' => false, 'primary_key' => 'id'])
            ->addColumn('id', 'biginteger', ['identity' => true])
            ->addColumn('record_id', 'biginteger', ['null' => true])
            ->addColumn('name', 'string', ['null' => true])
            ->addColumn('memo', 'string', ['null' => true])
            ->addColumn('created_at', 'timestamp', ['null' => true])
            ->addColumn('updated_at', 'timestamp', ['null' => true])
            ->addColumn('deleted_at', 'timestamp', ['null' => true])
            ->addForeignKey('record_id', 'service_records', 'id', ['delete' => 'SET NULL', 'update' => 'CASCADE'])
            ->save();

        $this->execute("
            CREATE TRIGGER `" . $table . "_created_at` BEFORE INSERT ON `" . $table . "` FOR EACH ROW BEGIN SET NEW.created_at = CURRENT_TIMESTAMP; END;
            CREATE TRIGGER `" . $table . "_updated_at` BEFORE UPDATE ON `" . $table . "` FOR EACH ROW BEGIN SET NEW.updated_at = CURRENT_TIMESTAMP; END;
        ");
    }
}
