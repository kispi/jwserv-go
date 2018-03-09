<?php

use Phinx\Migration\AbstractMigration;

class AddTriggersOnTime extends AbstractMigration
{
    protected $tables_to_create = [
        "congregations",
        "service_records",
        "users"
    ];

    public function up()
    {
        foreach ($this->tables_to_create as $table) {
            $table = $this->table($table);
            $table->addColumn('created_at', 'timestamp', array('null' => true))
                  ->addColumn('updated_at', 'timestamp', array('null' => true))
                  ->addColumn('deleted_at', 'timestamp', array('null' => true))
                  ->save();
        }

        foreach ($this->tables_to_create as $table) {
            $this->execute("
                CREATE TRIGGER `" . $table . "_created_at` BEFORE INSERT ON `" . $table . "` FOR EACH ROW BEGIN SET NEW.created_at = CURRENT_TIMESTAMP; END;
                CREATE TRIGGER `" . $table . "_updated_at` BEFORE UPDATE ON `" . $table . "` FOR EACH ROW BEGIN SET NEW.updated_at = CURRENT_TIMESTAMP; END;
            ");
        }
    }
}
