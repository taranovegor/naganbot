<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20200425145606 extends AbstractMigration
{
    public function getDescription() : string
    {
        return 'rr-17';
    }

    public function up(Schema $schema) : void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE gunslinger DROP FOREIGN KEY FK_A99081A0E48FD905');
        $this->addSql('ALTER TABLE shot_himself DROP FOREIGN KEY FK_6FBCC7CB6DAEA2B7');
        $this->addSql('CREATE TABLE game_game (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', chat_id BIGINT NOT NULL, owner_user_id INT NOT NULL, created_at DATETIME NOT NULL, played_out_at DATETIME DEFAULT NULL, INDEX IDX_C83E5DA01A9A7125 (chat_id), INDEX IDX_C83E5DA02B18554A (owner_user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE game_gunslinger (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', game_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', user_id INT NOT NULL, joined_at DATETIME NOT NULL, shot_himself TINYINT(1) DEFAULT \'0\' NOT NULL, INDEX IDX_474C464DE48FD905 (game_id), INDEX IDX_474C464DA76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('ALTER TABLE game_game ADD CONSTRAINT FK_C83E5DA01A9A7125 FOREIGN KEY (chat_id) REFERENCES telegram_chat (id)');
        $this->addSql('ALTER TABLE game_game ADD CONSTRAINT FK_C83E5DA02B18554A FOREIGN KEY (owner_user_id) REFERENCES telegram_user (id)');
        $this->addSql('ALTER TABLE game_gunslinger ADD CONSTRAINT FK_474C464DE48FD905 FOREIGN KEY (game_id) REFERENCES game_game (id)');
        $this->addSql('ALTER TABLE game_gunslinger ADD CONSTRAINT FK_474C464DA76ED395 FOREIGN KEY (user_id) REFERENCES telegram_user (id)');
        $this->addSql('INSERT INTO game_game(id, chat_id, owner_user_id, created_at, played_out_at) SELECT id, chat_id, owner_user_id, created_at, played_out_at FROM game');
        $this->addSql('INSERT INTO game_gunslinger(id, game_id, user_id, joined_at) SELECT id, game_id, user_id, joined_at FROM gunslinger');
        $this->addSql('UPDATE game_gunslinger gg INNER JOIN shot_himself sh on gg.id = sh.gunslinger_id SET gg.shot_himself = 1 WHERE gg.id = sh.gunslinger_id');
        $this->addSql('DROP TABLE game');
        $this->addSql('DROP TABLE gunslinger');
        $this->addSql('DROP TABLE shot_himself');
        $this->addSql('ALTER TABLE game_gunslinger CHANGE shot_himself shot_himself TINYINT(1) NOT NULL');
    }

    public function down(Schema $schema) : void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE game_gunslinger DROP FOREIGN KEY FK_474C464DE48FD905');
        $this->addSql('CREATE TABLE game (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', chat_id BIGINT NOT NULL, owner_user_id INT NOT NULL, created_at DATETIME NOT NULL, played_out_at DATETIME DEFAULT NULL, INDEX IDX_232B318C1A9A7125 (chat_id), INDEX IDX_232B318C2B18554A (owner_user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE `utf8_unicode_ci` ENGINE = InnoDB COMMENT = \'\' ');
        $this->addSql('CREATE TABLE gunslinger (id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', game_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', user_id INT NOT NULL, joined_at DATETIME NOT NULL, INDEX IDX_A99081A0E48FD905 (game_id), INDEX IDX_A99081A0A76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE `utf8_unicode_ci` ENGINE = InnoDB COMMENT = \'\' ');
        $this->addSql('CREATE TABLE shot_himself (gunslinger_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\', PRIMARY KEY(gunslinger_id)) DEFAULT CHARACTER SET utf8 COLLATE `utf8_unicode_ci` ENGINE = InnoDB COMMENT = \'\' ');
        $this->addSql('ALTER TABLE game ADD CONSTRAINT FK_232B318C1A9A7125 FOREIGN KEY (chat_id) REFERENCES telegram_chat (id)');
        $this->addSql('ALTER TABLE game ADD CONSTRAINT FK_232B318C2B18554A FOREIGN KEY (owner_user_id) REFERENCES telegram_user (id)');
        $this->addSql('ALTER TABLE gunslinger ADD CONSTRAINT FK_A99081A0A76ED395 FOREIGN KEY (user_id) REFERENCES telegram_user (id)');
        $this->addSql('ALTER TABLE gunslinger ADD CONSTRAINT FK_A99081A0E48FD905 FOREIGN KEY (game_id) REFERENCES game (id)');
        $this->addSql('ALTER TABLE shot_himself ADD CONSTRAINT FK_6FBCC7CB6DAEA2B7 FOREIGN KEY (gunslinger_id) REFERENCES gunslinger (id)');
        $this->addSql('INSERT INTO game(id, chat_id, owner_user_id, created_at, played_out_at) SELECT id, chat_id, owner_user_id, created_at, played_out_at FROM game_game');
        $this->addSql('INSERT INTO gunslinger(id, game_id, user_id, joined_at) SELECT id, game_id, user_id, joined_at FROM game_gunslinger');
        $this->addSql('INSERT INTO shot_himself(gunslinger_id) SELECT id FROM game_gunslinger WHERE shot_himself = 1');
        $this->addSql('DROP TABLE game_game');
        $this->addSql('DROP TABLE game_gunslinger');
    }
}
