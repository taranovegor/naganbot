<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20200814125753 extends AbstractMigration
{
    public function getDescription() : string
    {
        return '';
    }

    public function up(Schema $schema) : void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE game_game CHANGE played_out_at played_out_at DATETIME DEFAULT NULL');
        $this->addSql('ALTER TABLE telegram_chat CHANGE title title VARCHAR(255) DEFAULT NULL, CHANGE first_name first_name VARCHAR(255) DEFAULT NULL, CHANGE last_name last_name VARCHAR(255) DEFAULT NULL, CHANGE invite_link invite_link VARCHAR(255) DEFAULT NULL, CHANGE type type_type VARCHAR(255) NOT NULL');
        $this->addSql('ALTER TABLE telegram_user CHANGE username username_username VARCHAR(255) DEFAULT NULL, CHANGE language_code language_code_code VARCHAR(255) DEFAULT NULL, CHANGE last_name last_name VARCHAR(255) DEFAULT NULL');
    }

    public function down(Schema $schema) : void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE game_game CHANGE played_out_at played_out_at DATETIME DEFAULT NULL');
        $this->addSql('ALTER TABLE telegram_chat CHANGE title title VARCHAR(255) DEFAULT NULL, CHANGE first_name first_name VARCHAR(255) DEFAULT NULL, CHANGE last_name last_name VARCHAR(255) DEFAULT NULL, CHANGE invite_link invite_link VARCHAR(255) DEFAULT NULL, CHANGE type_type type VARCHAR(255) NOT NULL');
        $this->addSql('ALTER TABLE telegram_user CHANGE username_username username VARCHAR(255) DEFAULT NULL, CHANGE language_code_code language_code VARCHAR(255) DEFAULT NULL, CHANGE last_name last_name VARCHAR(255) DEFAULT NULL');
    }
}
