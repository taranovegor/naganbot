<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20200106200233 extends AbstractMigration
{
    public function getDescription() : string
    {
        return '';
    }

    public function up(Schema $schema) : void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('CREATE TABLE game (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', chat_id BIGINT NOT NULL, owner_user_id INT NOT NULL, created_at DATETIME NOT NULL, played_out_at DATETIME DEFAULT NULL, INDEX IDX_232B318C1A9A7125 (chat_id), INDEX IDX_232B318C2B18554A (owner_user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE shot_himself (gunslinger_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', PRIMARY KEY(gunslinger_id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE telegram_user (id INT NOT NULL, is_bot TINYINT(1) NOT NULL, first_name VARCHAR(255) NOT NULL, last_name VARCHAR(255) DEFAULT NULL, username VARCHAR(255) DEFAULT NULL, language_code VARCHAR(255) DEFAULT NULL, PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE telegram_chat (id BIGINT NOT NULL, type VARCHAR(255) NOT NULL, title VARCHAR(255) DEFAULT NULL, first_name VARCHAR(255) DEFAULT NULL, last_name VARCHAR(255) DEFAULT NULL, invite_link VARCHAR(255) DEFAULT NULL, PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE gunslinger (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', game_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', user_id INT NOT NULL, joined_at DATETIME NOT NULL, INDEX IDX_A99081A0E48FD905 (game_id), INDEX IDX_A99081A0A76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('ALTER TABLE game ADD CONSTRAINT FK_232B318C1A9A7125 FOREIGN KEY (chat_id) REFERENCES telegram_chat (id)');
        $this->addSql('ALTER TABLE game ADD CONSTRAINT FK_232B318C2B18554A FOREIGN KEY (owner_user_id) REFERENCES telegram_user (id)');
        $this->addSql('ALTER TABLE shot_himself ADD CONSTRAINT FK_6FBCC7CB6DAEA2B7 FOREIGN KEY (gunslinger_id) REFERENCES gunslinger (id)');
        $this->addSql('ALTER TABLE gunslinger ADD CONSTRAINT FK_A99081A0E48FD905 FOREIGN KEY (game_id) REFERENCES game (id)');
        $this->addSql('ALTER TABLE gunslinger ADD CONSTRAINT FK_A99081A0A76ED395 FOREIGN KEY (user_id) REFERENCES telegram_user (id)');
    }

    public function down(Schema $schema) : void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE gunslinger DROP FOREIGN KEY FK_A99081A0E48FD905');
        $this->addSql('ALTER TABLE game DROP FOREIGN KEY FK_232B318C2B18554A');
        $this->addSql('ALTER TABLE gunslinger DROP FOREIGN KEY FK_A99081A0A76ED395');
        $this->addSql('ALTER TABLE game DROP FOREIGN KEY FK_232B318C1A9A7125');
        $this->addSql('ALTER TABLE shot_himself DROP FOREIGN KEY FK_6FBCC7CB6DAEA2B7');
        $this->addSql('DROP TABLE game');
        $this->addSql('DROP TABLE shot_himself');
        $this->addSql('DROP TABLE telegram_user');
        $this->addSql('DROP TABLE telegram_chat');
        $this->addSql('DROP TABLE gunslinger');
    }
}
