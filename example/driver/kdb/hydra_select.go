package kdb

import (
	"database/sql"
	"fmt"
)

func TestHydraSelect(op *sql.DB) {
	tableSql := "CREATE TABLE IF NOT EXISTS `hydra_client` (" +
		"`id` VARCHAR(255) NOT NULL," +
		"`client_name` text NOT NULL," +
		"`client_secret` text NOT NULL," +
		"`scope` text NOT NULL," +
		"`owner` text NOT NULL," +
		"`policy_uri` text NOT NULL," +
		"`tos_uri` text NOT NULL," +
		"`client_uri` text NOT NULL," +
		"`logo_uri` text NOT NULL," +
		"`client_secret_expires_at` INT NOT NULL DEFAULT '0'," +
		"`sector_identifier_uri` text NOT NULL," +
		"`jwks` text NOT NULL," +
		"`jwks_uri` text NOT NULL," +
		"`token_endpoint_auth_method` VARCHAR(25) NOT NULL DEFAULT ''," +
		"`request_object_signing_alg` VARCHAR(10) NOT NULL DEFAULT ''," +
		"`userinfo_signed_response_alg` VARCHAR(10) NOT NULL DEFAULT ''," +
		"`subject_type` VARCHAR(15) NOT NULL DEFAULT ''," +
		"`pk_deprecated` INT  DEFAULT NULL," +
		"`created_at` TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`updated_at` TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`frontchannel_logout_uri` text NOT NULL," +
		"`frontchannel_logout_session_required` TINYINT(1) NOT NULL DEFAULT '0'," +
		"`backchannel_logout_uri` text NOT NULL," +
		"`backchannel_logout_session_required` TINYINT(1) NOT NULL DEFAULT '0'," +
		"`metadata` text NOT NULL," +
		"`token_endpoint_auth_signing_alg` VARCHAR(10) NOT NULL DEFAULT ''," +
		"`authorization_code_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`authorization_code_grant_id_token_lifespan` BIGINT DEFAULT NULL," +
		"`authorization_code_grant_refresh_token_lifespan` BIGINT DEFAULT NULL," +
		"`client_credentials_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`implicit_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`implicit_grant_id_token_lifespan` BIGINT DEFAULT NULL," +
		"`jwt_bearer_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`password_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`password_grant_refresh_token_lifespan` BIGINT DEFAULT NULL," +
		"`refresh_token_grant_id_token_lifespan` BIGINT DEFAULT NULL," +
		"`refresh_token_grant_access_token_lifespan` BIGINT DEFAULT NULL," +
		"`refresh_token_grant_refresh_token_lifespan` BIGINT DEFAULT NULL," +
		"`pk` VARCHAR(36) NOT NULL," +
		"`registration_access_token_signature` VARCHAR(128) NOT NULL DEFAULT ''," +
		"`nid` VARCHAR(36) NOT NULL," +
		"`redirect_uris` json NOT NULL," +
		"`grant_types` json NOT NULL," +
		"`response_types` json NOT NULL," +
		"`audience` json NOT NULL," +
		"`allowed_cors_origins` json NOT NULL," +
		"`contacts` json NOT NULL," +
		"`request_uris` json NOT NULL," +
		"`post_logout_redirect_uris` json NOT NULL," +
		"`access_token_strategy` VARCHAR(10) NOT NULL DEFAULT ''," +
		"`skip_consent` TINYINT(1) NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`pk`));"

	_, err := op.Exec(tableSql)
	if err != nil {
		fmt.Println(err)
		return
	}

	args := []any{1, "cb51aed3-f8f1-4c7f-aa3b-495a84219e65"}
	row := op.QueryRow("SELECT `hydra_client`.`access_token_strategy`, `hydra_client`.`allowed_cors_origins`, "+
		"`hydra_client`.`audience`, `hydra_client`.`authorization_code_grant_access_token_lifespan`, "+
		"`hydra_client`.`authorization_code_grant_id_token_lifespan`, "+
		"`hydra_client`.`authorization_code_grant_refresh_token_lifespan`, "+
		"`hydra_client`.`backchannel_logout_session_required`, `hydra_client`.`backchannel_logout_uri`, "+
		"`hydra_client`.`client_credentials_grant_access_token_lifespan`, `hydra_client`.`client_name`, "+
		"`hydra_client`.`client_secret_expires_at`, `hydra_client`.`client_secret`, "+
		"`hydra_client`.`client_uri`, `hydra_client`.`contacts`, `hydra_client`.`created_at`, "+
		"`hydra_client`.`frontchannel_logout_session_required`, `hydra_client`.`frontchannel_logout_uri`, "+
		"`hydra_client`.`grant_types`, `hydra_client`.`id`, `hydra_client`.`implicit_grant_access_token_lifespan`, "+
		"`hydra_client`.`implicit_grant_id_token_lifespan`, `hydra_client`.`jwks_uri`, `hydra_client`.`jwks`, "+
		"`hydra_client`.`jwt_bearer_grant_access_token_lifespan`, `hydra_client`.`logo_uri`, "+
		"`hydra_client`.`metadata`, `hydra_client`.`nid`, `hydra_client`.`owner`, "+
		"`hydra_client`.`password_grant_access_token_lifespan`, "+
		"`hydra_client`.`password_grant_refresh_token_lifespan`, `hydra_client`.`pk_deprecated`, "+
		"`hydra_client`.`pk`, `hydra_client`.`policy_uri`, `hydra_client`.`post_logout_redirect_uris`, "+
		"`hydra_client`.`redirect_uris`, `hydra_client`.`refresh_token_grant_access_token_lifespan`, "+
		"`hydra_client`.`refresh_token_grant_id_token_lifespan`, "+
		"`hydra_client`.`refresh_token_grant_refresh_token_lifespan`, "+
		"`hydra_client`.`registration_access_token_signature`, `hydra_client`.`request_object_signing_alg`, "+
		"`hydra_client`.`request_uris`, `hydra_client`.`response_types`, `hydra_client`.`scope`, "+
		"`hydra_client`.`sector_identifier_uri`, `hydra_client`.`skip_consent`, `hydra_client`.`subject_type`, "+
		"`hydra_client`.`token_endpoint_auth_method`, `hydra_client`.`token_endpoint_auth_signing_alg`, "+
		"`hydra_client`.`tos_uri`, `hydra_client`.`updated_at`, `hydra_client`.`userinfo_signed_response_alg` "+
		"FROM hydra_client AS hydra_client WHERE nid = ? AND id = ?", args...)
	err = row.Scan()
	if err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}

	fmt.Println("success")
}
