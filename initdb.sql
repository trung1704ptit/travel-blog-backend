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
  CONSTRAINT `category_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `category` (`id`) ON DELETE SET NULL
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
  `category_id` char(36) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `composite` (`article_id`,`category_id`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `article_category_ibfk_1` FOREIGN KEY (`article_id`) REFERENCES `article` (`id`) ON DELETE CASCADE,
  CONSTRAINT `article_category_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE SET NULL
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
-- Main Categories (Root Level)
('10000000-0000-0000-0000-000000000001','UAE Destinations','uae-destinations','Explore cities and attractions across the UAE',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('20000000-0000-0000-0000-000000000001','Hotels & Resorts','hotels-resorts','Luxury and budget accommodations in UAE',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('30000000-0000-0000-0000-000000000001','Entertainment','entertainment','Fun activities and entertainment venues',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('40000000-0000-0000-0000-000000000001','Gaming & Casinos','gaming-casinos','Gaming venues and entertainment complexes',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('50000000-0000-0000-0000-000000000001','Dining & Nightlife','dining-nightlife','Restaurants, cafes, and nightlife spots',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('60000000-0000-0000-0000-000000000001','Travel Tips','travel-tips','Essential UAE travel information',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('70000000-0000-0000-0000-000000000001','Activities & Adventures','activities-adventures','Outdoor and indoor activities',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),
('80000000-0000-0000-0000-000000000001','Shopping','shopping','Malls, souks, and shopping destinations',NULL,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- UAE Destinations Subcategories
('11000000-0000-0000-0000-000000000001','Dubai','dubai','The city of superlatives','https://example.com/images/dubai.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000002','Abu Dhabi','abu-dhabi','UAE capital city','https://example.com/images/abudhabi.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000003','Sharjah','sharjah','Cultural capital of UAE','https://example.com/images/sharjah.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000004','Ajman','ajman','Peaceful coastal emirate','https://example.com/images/ajman.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000005','Ras Al Khaimah','ras-al-khaimah','Mountain and beach paradise','https://example.com/images/rak.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000006','Fujairah','fujairah','East coast beauty','https://example.com/images/fujairah.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11000000-0000-0000-0000-000000000007','Umm Al Quwain','umm-al-quwain','Hidden gem emirate','https://example.com/images/uaq.jpg','10000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Dubai Subcategories (Level 2)
('11100000-0000-0000-0000-000000000001','Downtown Dubai','downtown-dubai','Burj Khalifa and Dubai Mall area','https://example.com/images/downtown-dubai.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000002','Dubai Marina','dubai-marina','Waterfront living and dining','https://example.com/images/marina.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000003','Palm Jumeirah','palm-jumeirah','Iconic palm-shaped island','https://example.com/images/palm.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000004','JBR Beach','jbr-beach','Jumeirah Beach Residence','https://example.com/images/jbr.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000005','Old Dubai','old-dubai','Historic Al Fahidi and Creek','https://example.com/images/old-dubai.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000006','Dubai Creek','dubai-creek','Traditional waterway area','https://example.com/images/creek.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11100000-0000-0000-0000-000000000007','Business Bay','business-bay','Modern business district','https://example.com/images/business-bay.jpg','11000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Abu Dhabi Subcategories (Level 2)
('11200000-0000-0000-0000-000000000001','Yas Island','yas-island','Entertainment and leisure hub','https://example.com/images/yas.jpg','11000000-0000-0000-0000-000000000002','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11200000-0000-0000-0000-000000000002','Saadiyat Island','saadiyat-island','Cultural island destination','https://example.com/images/saadiyat.jpg','11000000-0000-0000-0000-000000000002','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11200000-0000-0000-0000-000000000003','Corniche','corniche','Waterfront promenade','https://example.com/images/corniche.jpg','11000000-0000-0000-0000-000000000002','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('11200000-0000-0000-0000-000000000004','Al Ain','al-ain','Garden city oasis','https://example.com/images/alain.jpg','11000000-0000-0000-0000-000000000002','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Hotels & Resorts Subcategories
('21000000-0000-0000-0000-000000000001','7-Star Hotels','7-star-hotels','Ultra-luxury properties','https://example.com/images/7star.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('21000000-0000-0000-0000-000000000002','5-Star Resorts','5-star-resorts','Premium beachfront resorts','https://example.com/images/5star.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('21000000-0000-0000-0000-000000000003','Beach Resorts','beach-resorts','Coastal accommodations','https://example.com/images/beach-resort.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('21000000-0000-0000-0000-000000000004','City Hotels','city-hotels','Urban accommodation','https://example.com/images/city-hotel.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('21000000-0000-0000-0000-000000000005','Budget Hotels','budget-hotels','Affordable stays','https://example.com/images/budget.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('21000000-0000-0000-0000-000000000006','Apartments & Rentals','apartments-rentals','Serviced apartments','https://example.com/images/apartment.jpg','20000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Entertainment Subcategories
('31000000-0000-0000-0000-000000000001','Theme Parks','theme-parks','Adventure and theme parks','https://example.com/images/theme-park.jpg','30000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('31000000-0000-0000-0000-000000000002','Water Parks','water-parks','Aquatic fun and slides','https://example.com/images/waterpark.jpg','30000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('31000000-0000-0000-0000-000000000003','Indoor Entertainment','indoor-entertainment','Indoor fun centers','https://example.com/images/indoor.jpg','30000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('31000000-0000-0000-0000-000000000004','Shows & Performances','shows-performances','Live entertainment','https://example.com/images/shows.jpg','30000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('31000000-0000-0000-0000-000000000005','Cinemas & Theaters','cinemas-theaters','Movie theaters and venues','https://example.com/images/cinema.jpg','30000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Gaming & Casinos Subcategories
('41000000-0000-0000-0000-000000000001','Resort Casinos','resort-casinos','Integrated casino resorts','https://example.com/images/casino.jpg','40000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('41000000-0000-0000-0000-000000000002','Gaming Lounges','gaming-lounges','Premium gaming venues','https://example.com/images/gaming.jpg','40000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('41000000-0000-0000-0000-000000000003','eSports Venues','esports-venues','Competitive gaming centers','https://example.com/images/esports.jpg','40000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('41000000-0000-0000-0000-000000000004','Entertainment Complexes','entertainment-complexes','Mixed entertainment venues','https://example.com/images/complex.jpg','40000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Dining & Nightlife Subcategories
('51000000-0000-0000-0000-000000000001','Fine Dining','fine-dining','Michelin-star restaurants','https://example.com/images/finedining.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000002','Arabic Cuisine','arabic-cuisine','Traditional Emirati food','https://example.com/images/arabic.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000003','International Cuisine','international-cuisine','Global dining options','https://example.com/images/international.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000004','Cafes & Desserts','cafes-desserts','Coffee shops and sweet treats','https://example.com/images/cafe.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000005','Rooftop Bars','rooftop-bars','Sky-high dining and drinks','https://example.com/images/rooftop.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000006','Beach Clubs','beach-clubs','Seaside lounges','https://example.com/images/beachclub.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('51000000-0000-0000-0000-000000000007','Nightclubs','nightclubs','Late-night entertainment','https://example.com/images/nightclub.jpg','50000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Travel Tips Subcategories
('61000000-0000-0000-0000-000000000001','Visa & Immigration','visa-immigration','Entry requirements','https://example.com/images/visa.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('61000000-0000-0000-0000-000000000002','Transportation','transportation','Getting around UAE','https://example.com/images/transport.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('61000000-0000-0000-0000-000000000003','Weather & Best Time','weather-best-time','Climate and seasons','https://example.com/images/weather.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('61000000-0000-0000-0000-000000000004','Culture & Customs','culture-customs','Local traditions','https://example.com/images/culture.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('61000000-0000-0000-0000-000000000005','Money & Currency','money-currency','AED and payments','https://example.com/images/money.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('61000000-0000-0000-0000-000000000006','Safety Tips','safety-tips','Travel safety advice','https://example.com/images/safety.jpg','60000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Activities & Adventures Subcategories
('71000000-0000-0000-0000-000000000001','Desert Safari','desert-safari','Dune bashing and camping','https://example.com/images/desert.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('71000000-0000-0000-0000-000000000002','Water Sports','water-sports','Jet skiing and diving','https://example.com/images/watersports.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('71000000-0000-0000-0000-000000000003','Skydiving','skydiving','Tandem and solo jumps','https://example.com/images/skydive.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('71000000-0000-0000-0000-000000000004','Yacht Cruises','yacht-cruises','Luxury boat tours','https://example.com/images/yacht.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('71000000-0000-0000-0000-000000000005','Golf','golf','World-class golf courses','https://example.com/images/golf.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('71000000-0000-0000-0000-000000000006','Spa & Wellness','spa-wellness','Relaxation and treatments','https://example.com/images/spa.jpg','70000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),

-- Shopping Subcategories
('81000000-0000-0000-0000-000000000001','Luxury Malls','luxury-malls','Premium shopping centers','https://example.com/images/mall.jpg','80000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('81000000-0000-0000-0000-000000000002','Traditional Souks','traditional-souks','Gold, spice, and textile markets','https://example.com/images/souk.jpg','80000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('81000000-0000-0000-0000-000000000003','Outlet Shopping','outlet-shopping','Discounted brand outlets','https://example.com/images/outlet.jpg','80000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00'),
('81000000-0000-0000-0000-000000000004','Shopping Festivals','shopping-festivals','Dubai Shopping Festival','https://example.com/images/festival.jpg','80000000-0000-0000-0000-000000000001','2024-01-01 00:00:00','2024-01-01 00:00:00');

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
