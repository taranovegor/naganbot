<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20200829125253 extends AbstractMigration
{
    public function getDescription() : string
    {
        return 'ng-11';
    }

    public function up(Schema $schema) : void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('SET FOREIGN_KEY_CHECKS = 0;');
        $this->addSql('DROP FUNCTION IF EXISTS ouuid_to_uuid;');
        $this->addSql('CREATE FUNCTION ouuid_to_uuid(uuid BINARY(16)) RETURNS VARCHAR(36) RETURN LOWER(CONCAT(SUBSTR(HEX(uuid), 9, 8), \'-\', SUBSTR(HEX(uuid), 5, 4), \'-\', SUBSTR(HEX(uuid), 1, 4), \'-\', SUBSTR(HEX(uuid), 17,4), \'-\', SUBSTR(HEX(uuid), 21, 12 )));');

        $this->addSql('ALTER TABLE game_gunslinger CHANGE id `id_bin` BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\';');
        $this->addSql('ALTER TABLE game_gunslinger ADD COLUMN id CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\' FIRST;');
        $this->addSql('UPDATE game_gunslinger SET id = ouuid_to_uuid(id_bin) WHERE 1;');
        $this->addSql('ALTER TABLE game_gunslinger DROP COLUMN id_bin;');

        $this->addSql('ALTER TABLE game_gunslinger CHANGE game_id `game_id_bin` BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\';');
        $this->addSql('ALTER TABLE game_gunslinger ADD COLUMN game_id CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\' AFTER id;');
        $this->addSql('UPDATE game_gunslinger SET game_id = ouuid_to_uuid(game_id_bin) WHERE 1;');
        $this->addSql('ALTER TABLE game_gunslinger DROP FOREIGN KEY FK_474C464DE48FD905;');
        $this->addSql('ALTER TABLE game_gunslinger DROP INDEX IDX_474C464DE48FD905;');
        $this->addSql('ALTER TABLE game_gunslinger DROP COLUMN game_id_bin;');
        $this->addSql('ALTER TABLE game_gunslinger ADD PRIMARY KEY (id);');

        $this->addSql('ALTER TABLE game_game CHANGE id `id_bin` BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\';');
        $this->addSql('ALTER TABLE game_game ADD COLUMN id CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\' FIRST;');
        $this->addSql('UPDATE game_game SET id = ouuid_to_uuid(id_bin) WHERE 1;');
        $this->addSql('ALTER TABLE game_game DROP COLUMN id_bin;');
        $this->addSql('ALTER TABLE game_game ADD PRIMARY KEY (id);');

        $this->addSql('CREATE INDEX IDX_474C464DE48FD905 ON game_gunslinger (game_id);');
        $this->addSql('ALTER TABLE game_gunslinger ADD CONSTRAINT FK_474C464DE48FD905 FOREIGN KEY (game_id) REFERENCES game_game (id);');

        $this->addSql('DROP FUNCTION IF EXISTS ouuid_to_uuid;');
        $this->addSql('SET FOREIGN_KEY_CHECKS = 1;');

        $this->addSql('ALTER TABLE game_gunslinger CHANGE game_id game_id CHAR(36) NOT NULL COMMENT \'(DC2Type:uuid)\';');
        $this->addSql('ALTER TABLE game_game CHANGE played_out_at played_at DATETIME DEFAULT NULL;');

    }

    public function down(Schema $schema) : void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('SET FOREIGN_KEY_CHECKS = 0;');
        $this->addSql('DROP FUNCTION IF EXISTS uuid_to_ouuid;');
        $this->addSql('CREATE FUNCTION `uuid_to_ouuid`(uuid VARCHAR(36)) RETURNS BINARY(16) DETERMINISTIC RETURN UNHEX(CONCAT(SUBSTR(uuid, 15, 4), SUBSTR(uuid, 10, 4), SUBSTR(uuid, 1, 8), SUBSTR(uuid, 20, 4), SUBSTR(uuid, 25, 12)));');

        $this->addSql('ALTER TABLE game_gunslinger CHANGE id `id_str` CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\';');
        $this->addSql('ALTER TABLE game_gunslinger ADD COLUMN id BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\' FIRST;');
        $this->addSql('UPDATE game_gunslinger SET id = uuid_to_ouuid(id_str) WHERE 1;');
        $this->addSql('ALTER TABLE game_gunslinger DROP COLUMN id_str;');

        $this->addSql('ALTER TABLE game_gunslinger CHANGE game_id `game_id_str` CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\';');
        $this->addSql('ALTER TABLE game_gunslinger ADD COLUMN game_id BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\' AFTER id;');
        $this->addSql('UPDATE game_gunslinger SET game_id = uuid_to_ouuid(game_id_str) WHERE 1;');
        $this->addSql('ALTER TABLE game_gunslinger DROP FOREIGN KEY FK_474C464DE48FD905;');
        $this->addSql('ALTER TABLE game_gunslinger DROP INDEX IDX_474C464DE48FD905;');
        $this->addSql('ALTER TABLE game_gunslinger DROP COLUMN game_id_str;');
        $this->addSql('ALTER TABLE game_gunslinger ADD PRIMARY KEY (id);');

        $this->addSql('ALTER TABLE game_game CHANGE id `id_str` CHAR(36) NULL DEFAULT NULL COMMENT \'(DC2Type:uuid)\';');
        $this->addSql('ALTER TABLE game_game ADD COLUMN id BINARY(16) NULL COMMENT \'(DC2Type:uuid_binary)\' FIRST;');
        $this->addSql('UPDATE game_game SET id = uuid_to_ouuid(id_str) WHERE 1;');
        $this->addSql('ALTER TABLE game_game DROP COLUMN id_str;');
        $this->addSql('ALTER TABLE game_game ADD PRIMARY KEY (id);');

        $this->addSql('CREATE INDEX IDX_474C464DE48FD905 ON game_gunslinger (game_id);');
        $this->addSql('ALTER TABLE game_gunslinger ADD CONSTRAINT FK_474C464DE48FD905 FOREIGN KEY (game_id) REFERENCES game_game (id);');

        $this->addSql('DROP FUNCTION IF EXISTS uuid_to_ouuid;');
        $this->addSql('SET FOREIGN_KEY_CHECKS = 1;');

        $this->addSql('ALTER TABLE game_game CHANGE played_at played_out_at DATETIME DEFAULT NULL;');
        $this->addSql('ALTER TABLE game_gunslinger CHANGE game_id game_id BINARY(16) NOT NULL COMMENT \'(DC2Type:uuid_binary)\';');
    }
}
