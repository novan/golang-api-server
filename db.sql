-- Create syntax for TABLE 'users'
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_type` enum('ADMIN','USER') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'USER',
  `email` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  `mobile` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login` datetime DEFAULT NULL,
  `is_active` bit(1) NOT NULL DEFAULT b'1',
  `created_by` int(11) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_by` int(11) NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_email_idx` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Create syntax for TABLE 'user_tokens'
CREATE TABLE `user_tokens` (
  `token` varchar(64) COLLATE utf8_unicode_ci NOT NULL,
  `jwt_id` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `expired_at` datetime DEFAULT NULL,
  `is_used` bit(1) NOT NULL DEFAULT b'0',
  `invalidated` bit(1) NOT NULL DEFAULT b'0',
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`token`),
  KEY `user_id_idx` (`user_id`),
  CONSTRAINT `user_tokens_ibfk1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;