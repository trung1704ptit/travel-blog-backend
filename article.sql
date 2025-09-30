-- Migration script to clean and recreate database with updated schema
-- Run this script to update your database structure

-- Drop existing database and recreate
DROP DATABASE IF EXISTS `article`;
CREATE DATABASE IF NOT EXISTS `article` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `article`;

-- Set SQL mode and timezone
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `author`
--
DROP TABLE IF EXISTS `author`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `author` (
  `id` char(36) NOT NULL,
  `name` varchar(200) COLLATE utf8_unicode_ci DEFAULT '""',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `category`
--
DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `category` (
  `id` char(36) NOT NULL,
  `name` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `slug` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `description` text COLLATE utf8_unicode_ci DEFAULT NULL,
  `image` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `parent_id` char(36) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `category_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `category` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `article`
--
DROP TABLE IF EXISTS `article`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `article` (
  `id` char(36) NOT NULL,
  `title` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `slug` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `content` longtext COLLATE utf8_unicode_ci NOT NULL,
  `thumbnail` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `image` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `short_description` text COLLATE utf8_unicode_ci DEFAULT NULL,
  `meta_description` text COLLATE utf8_unicode_ci DEFAULT NULL,
  `keywords` json DEFAULT NULL,
  `tags` json DEFAULT NULL,
  `reading_time_minutes` int DEFAULT 0,
  `views` int DEFAULT 0,
  `likes` int DEFAULT 0,
  `comments` int DEFAULT 0,
  `published` boolean DEFAULT false,
  `published_at` datetime DEFAULT NULL,
  `author_id` char(36) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `article_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `author` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `article_category`
--
DROP TABLE IF EXISTS `article_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `article_category` (
  `id` char(36) NOT NULL,
  `article_id` char(36) NOT NULL,
  `category_id` char(36) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `composite` (`article_id`,`category_id`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `article_category_ibfk_1` FOREIGN KEY (`article_id`) REFERENCES `article` (`id`) ON DELETE CASCADE,
  CONSTRAINT `article_category_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Insert sample data
--

-- Insert authors
LOCK TABLES `author` WRITE;
/*!40000 ALTER TABLE `author` DISABLE KEYS */;
INSERT INTO `author` VALUES 
('550e8400-e29b-41d4-a716-446655440000','Iman Tumorang','2017-05-18 13:50:19','2017-05-18 13:50:19');
/*!40000 ALTER TABLE `author` ENABLE KEYS */;
UNLOCK TABLES;

-- Insert categories with nested structure
LOCK TABLES `category` WRITE;
/*!40000 ALTER TABLE `category` DISABLE KEYS */;
INSERT INTO `category` VALUES 
-- Root categories
('550e8400-e29b-41d4-a716-446655440001','Makanan','food','Kategori untuk semua hal tentang makanan','https://example.com/images/food-category.jpg',NULL,'2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440002','Kehidupan','life','Kategori untuk semua hal tentang kehidupan','https://example.com/images/life-category.jpg',NULL,'2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440003','Kasih Sayang','love','Kategori untuk semua hal tentang kasih sayang','https://example.com/images/love-category.jpg',NULL,'2017-05-18 13:50:19','2017-05-18 13:50:19'),
-- Sub-categories of Makanan
('550e8400-e29b-41d4-a716-446655440004','Masakan Indonesia','indonesian-food','Masakan tradisional Indonesia','https://example.com/images/indonesian-food.jpg','550e8400-e29b-41d4-a716-446655440001','2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440005','Masakan Barat','western-food','Masakan dari negara-negara barat','https://example.com/images/western-food.jpg','550e8400-e29b-41d4-a716-446655440001','2017-05-18 13:50:19','2017-05-18 13:50:19'),
-- Sub-categories of Masakan Indonesia
('550e8400-e29b-41d4-a716-446655440006','Nasi Goreng','nasi-goreng','Berbagai jenis nasi goreng','https://example.com/images/nasi-goreng.jpg','550e8400-e29b-41d4-a716-446655440004','2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440007','Soto','soto','Berbagai jenis soto','https://example.com/images/soto.jpg','550e8400-e29b-41d4-a716-446655440004','2017-05-18 13:50:19','2017-05-18 13:50:19');
/*!40000 ALTER TABLE `category` ENABLE KEYS */;
UNLOCK TABLES;

-- Insert articles
LOCK TABLES `article` WRITE;
/*!40000 ALTER TABLE `article` DISABLE KEYS */;
INSERT INTO `article` VALUES 
('550e8400-e29b-41d4-a716-446655440010','Makan Ayam','makan-ayam','<p>But I must explain to you how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system...</p>','https://example.com/thumb1.jpg','https://example.com/img1.jpg','A delicious article about eating chicken','Meta description for chicken article','["food", "chicken", "recipe"]','["cooking", "healthy"]',5,100,25,10,true,'2017-05-18 13:50:19','550e8400-e29b-41d4-a716-446655440000','2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440011','Makan Ikan','makan-ikan','<h1>Odio Mollis Turpis Dictumst</h1><p>Ut arcu tempor auctor pellentesque vitae lacinia potenti amet tellus sagittis molestie aliquam est mi facilisi amet...</p>','https://example.com/thumb2.jpg','https://example.com/img2.jpg','An article about eating fish','Meta description for fish article','["food", "fish", "seafood"]','["cooking", "seafood"]',7,150,30,15,true,'2017-05-18 13:50:19','550e8400-e29b-41d4-a716-446655440000','2017-05-18 13:50:19','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440012','Makan Sayur','makan-sayur','Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi id odio tortor. Pellentesque in efficitur velit...','https://example.com/thumb3.jpg','https://example.com/img3.jpg','A healthy article about eating vegetables','Meta description for vegetables article','["food", "vegetables", "healthy"]','["cooking", "healthy", "vegetarian"]',4,80,20,8,true,'2017-05-18 13:50:19','550e8400-e29b-41d4-a716-446655440000','2017-05-18 13:50:19','2017-05-18 13:50:19');
/*!40000 ALTER TABLE `article` ENABLE KEYS */;
UNLOCK TABLES;

-- Insert article-category relationships
LOCK TABLES `article_category` WRITE;
/*!40000 ALTER TABLE `article_category` DISABLE KEYS */;
INSERT INTO `article_category` VALUES 
('550e8400-e29b-41d4-a716-446655440020','550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440001','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440021','550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440002','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440022','550e8400-e29b-41d4-a716-446655440011','550e8400-e29b-41d4-a716-446655440001','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440023','550e8400-e29b-41d4-a716-446655440011','550e8400-e29b-41d4-a716-446655440002','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440024','550e8400-e29b-41d4-a716-446655440012','550e8400-e29b-41d4-a716-446655440001','2017-05-18 13:50:19'),
('550e8400-e29b-41d4-a716-446655440025','550e8400-e29b-41d4-a716-446655440012','550e8400-e29b-41d4-a716-446655440003','2017-05-18 13:50:19');
/*!40000 ALTER TABLE `article_category` ENABLE KEYS */;
UNLOCK TABLES;

-- Reset SQL mode and settings
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Migration completed successfully
SELECT 'Database migration completed successfully!' as status;
