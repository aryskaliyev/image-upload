-- forum.sql

-- PRAGMA foreign_keys = ON;

-- Create a table to store user accounts
CREATE TABLE IF NOT EXISTS useraccount (
	user_id INTEGER NOT NULL,
	username VARCHAR(30) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password BLOB NOT NULL,
	created DATETIME NOT NULL,
	UNIQUE (username),
	UNIQUE (email),
	PRIMARY KEY (user_id)
);

-- Create a table to store sessions
CREATE TABLE IF NOT EXISTS session (
	user_id INTEGER NOT NULL,
	uuid_token TEXT NOT NULL,
	created DATETIME NOT NULL,
	expires DATETIME NOT NULL,
	UNIQUE (user_id),
	FOREIGN KEY (user_id)
		REFERENCES useraccount(user_id)
		ON DELETE CASCADE,
	PRIMARY KEY (user_id)
);

-- Create a table to store categories
CREATE TABLE IF NOT EXISTS category (
	category_id INTEGER NOT NULL,
	name VARCHAR(20) NOT NULL,
	UNIQUE (name),
	PRIMARY KEY (category_id)	
);

-- Create a table to store posts
CREATE TABLE IF NOT EXISTS post (
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	title VARCHAR(75) NOT NULL,
	body VARCHAR(500) NOT NULL,
	image BLOB,
	created DATETIME NOT NULL,
	FOREIGN KEY (user_id)
		REFERENCES useraccount(user_id)
		ON DELETE CASCADE,
	PRIMARY KEY (post_id)
);

-- Create a table to store post-votes
CREATE TABLE IF NOT EXISTS post_vote (
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	vote INTEGER NOT NULL CHECK (vote == 1 OR vote == -1),
	FOREIGN KEY (user_id)
		REFERENCES useraccount(user_id)
		ON DELETE CASCADE,
	FOREIGN KEY (post_id)
		REFERENCES post(post_id)
		ON DELETE CASCADE,
	PRIMARY KEY (post_id, user_id)
);

-- Create a table to store post-categories
CREATE TABLE IF NOT EXISTS post_category (
	post_category_id INTEGER NOT NULL,
	post_id	INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	FOREIGN KEY (post_id)
		REFERENCES post(post_id)
		ON DELETE CASCADE,
	FOREIGN KEY (category_id)
		REFERENCES category(category_id)
		ON DELETE CASCADE,
	PRIMARY KEY (post_category_id)
);

-- Create a table to store comments
CREATE TABLE IF NOT EXISTS comment (
	comment_id INTEGER NOT NULL,
	body VARCHAR(75) NOT NULL,
	user_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
	created DATETIME NOT NULL,
	FOREIGN KEY (user_id)
		REFERENCES useraccount(user_id)
		ON DELETE CASCADE,
	FOREIGN KEY (post_id)
		REFERENCES post(post_id)
		ON DELETE CASCADE,
	PRIMARY KEY (comment_id)
);

-- Create a table to store comment-votes
CREATE TABLE IF NOT EXISTS comment_vote (
	comment_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	vote INTEGER NOT NULL CHECK (vote == 1 OR vote == -1),
	FOREIGN KEY (user_id)
		REFERENCES useraccount(user_id)
		ON DELETE CASCADE,
	FOREIGN KEY (comment_id)
		REFERENCES comment(comment_id)
		ON DELETE CASCADE,
	PRIMARY KEY (comment_id, user_id)
);

-- johndoepassword
INSERT INTO useraccount (username, email, hashed_password, created)
VALUES('john.doe', 'john.doe@email.com', '$2a$12$2VkdVil9bvDO3ZZ1DWWlMu0FLvKPbuj06bhxaXQqBv8JPbPzVsOQe', CURRENT_TIMESTAMP);

INSERT INTO useraccount (username, email, hashed_password, created)
VALUES('jane.smith', 'jane.smith@email.com', '$2a$12$2VkdVil9bvDO3ZZ1DWWlMu0FLvKPbuj06bhxaXQqBv8JPbPzVsOQe', CURRENT_TIMESTAMP);

INSERT INTO useraccount (username, email, hashed_password, created)
VALUES('jeremy.clarkson', 'jeremy.clarkson@email.com', '$2a$12$2VkdVil9bvDO3ZZ1DWWlMu0FLvKPbuj06bhxaXQqBv8JPbPzVsOQe', CURRENT_TIMESTAMP);

-- category data
INSERT INTO category (name) VALUES('cleancode');
INSERT INTO category (name) VALUES('golang');
INSERT INTO category (name) VALUES('backend');
INSERT INTO category (name) VALUES('html');
INSERT INTO category (name) VALUES('architecture');

-- post data
INSERT INTO post (user_id, title, body, created) VALUES(1, 'Proin a nisl quis dolor varius tristique ut sit amet metus', 'Nunc iaculis ipsum nec blandit porttitor. Etiam orci nibh, posuere lobortis pulvinar consequat, vulputate et est. Quisque mattis velit at nibh iaculis placerat. Morbi ut fermentum tellus, eget facilisis mi. Nullam et neque scelerisque, faucibus odio eget, vestibulum velit. Vivamus tempor nec erat vitae sodales. Pellentesque neque sapien, lacinia in malesuada vel, ultrices non orci. Quisque ligula sem, laoreet in euismod sit amet, semper maximus libero. Fusce at orci vel magna mollis imperdiet. Suspendisse potenti. Nunc urna neque, posuere eget suscipit eu, sollicitudin eu metus. Maecenas venenatis neque non lacus pellentesque, ornare egestas risus egestas. Praesent condimentum consequat lorem sit amet rutrum.

Nunc nisi quam, aliquam vel massa ac, blandit sagittis elit. Sed pharetra pharetra orci, at auctor odio sodales ac. Suspendisse tincidunt pulvinar rhoncus. Quisque et tempus felis. Ut eget lobortis ante. Morbi posuere velit ac justo bibendum semper. Morbi purus est, pellentesque vehicula venenatis at, mollis non libero. Duis vel lacus porta, egestas libero id, fermentum orci. Vivamus malesuada ipsum non libero pharetra porta. Pellentesque et lorem commodo, pellentesque dui nec, imperdiet enim. Nulla vel urna aliquet turpis condimentum efficitur eget at arcu. Aliquam congue laoreet diam, mollis consectetur ex ullamcorper quis.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(1, 'Praesent hendrerit odio vel malesuada laoreet', 'Praesent sollicitudin, mi quis efficitur mattis, quam est consequat dui, quis tempor nisl felis a metus. Mauris pulvinar fringilla nisi. Proin et varius justo. Donec sit amet nisl ac nisl rhoncus auctor. Mauris erat massa, egestas interdum erat nec, dapibus efficitur lectus. Nam a cursus nulla, ut maximus nisl. Nullam at ex vel eros consectetur suscipit. Pellentesque ullamcorper eget nunc nec tempor. Sed ante ligula, placerat sodales hendrerit vitae, tincidunt sit amet felis. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Donec vel posuere magna. Donec molestie suscipit risus, non laoreet sem vulputate eget. Phasellus porttitor bibendum congue. Curabitur scelerisque dui augue, quis aliquam dui molestie vel. Pellentesque facilisis interdum mauris vel hendrerit. In hac habitasse platea dictumst.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(2, 'Nunc tempus velit scelerisque sollicitudin porttitor', 'Morbi sit amet lobortis ante. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Suspendisse molestie nunc ac mollis varius. Curabitur a mauris viverra, eleifend ipsum nec, suscipit turpis. Ut vitae ullamcorper ligula. Sed ut arcu eu arcu dignissim pulvinar ultricies quis leo. Quisque interdum fringilla ex a pellentesque. Duis ac enim luctus, ornare tellus non, vulputate nulla. Mauris fermentum commodo velit et placerat.

Nam viverra interdum ex, sed dignissim lacus vestibulum eget. Pellentesque id massa egestas, bibendum ex a, consectetur nisi. Donec lobortis sapien id quam vestibulum ultrices. In hac habitasse platea dictumst. Morbi tempus tortor nec enim sagittis mattis. Integer eget arcu maximus, ultrices ante nec, vestibulum libero. Morbi lorem dui, tincidunt a convallis et, ultrices quis urna. Morbi mollis vitae enim in ultricies. Etiam vulputate venenatis nisi, sed accumsan erat scelerisque sed. Aliquam id eros sed dolor varius hendrerit. Curabitur viverra laoreet nibh, et maximus orci tempus at. Donec condimentum venenatis egestas. Ut feugiat quam at arcu varius, ac pulvinar diam consequat. ', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(2, 'Etiam et purus efficitur, varius elit nec, ornare lectus', 'Aenean nec condimentum nunc, eget ultrices enim. Donec sit amet malesuada libero. Maecenas purus augue, euismod a condimentum eget, sollicitudin a orci. Interdum et malesuada fames ac ante ipsum primis in faucibus. Ut diam ex, sagittis ac mattis non, convallis vitae nisi. Proin at fermentum magna, in congue sapien. Quisque et pharetra quam. Sed eget leo non purus porta fringilla. Quisque sit amet mi vitae orci facilisis rutrum. Integer turpis tortor, congue tempus tincidunt ut, tempus quis nisl.

Quisque eu pharetra odio. Aliquam lectus ex, feugiat eget mi sed, lobortis pretium tortor. Sed nec tortor ut urna volutpat lacinia. Proin laoreet, metus et porta ultrices, tortor tellus interdum mi, et ultrices nulla lacus id odio. Nam a cursus dolor. Fusce lectus lorem, tincidunt vitae lectus accumsan, ullamcorper luctus metus. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Morbi ultrices lorem sed sollicitudin gravida. Proin id odio in odio bibendum luctus id in justo. Curabitur faucibus non justo ac congue. Quisque purus nunc, gravida viverra tincidunt sit amet, consectetur at nisi.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(2, 'Donec interdum sem ac nisi congue dignissim', 'Etiam condimentum urna ac convallis sollicitudin. Donec nec nunc a tellus suscipit accumsan at eu sapien. Sed ut magna sed nisl mollis eleifend. Cras id augue placerat neque luctus laoreet. Donec vehicula malesuada urna, in mattis purus malesuada sit amet. Nunc volutpat massa non tempor venenatis. Maecenas erat urna, aliquam nec porttitor vel, porta at leo. Mauris egestas suscipit sapien nec placerat. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Morbi nec volutpat orci. Vivamus feugiat ullamcorper risus quis aliquam. Pellentesque gravida eros vel lectus malesuada lobortis. Praesent bibendum libero sit amet tempus faucibus. Pellentesque aliquam nunc et rutrum sodales. Phasellus finibus cursus felis eu elementum.' , CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(3, 'Aliquam hendrerit arcu eu fermentum scelerisque', 'Sed consectetur pellentesque vestibulum. Praesent convallis nisl eu ullamcorper fringilla. Aliquam egestas ut lorem non dignissim. Maecenas hendrerit non tellus sed vehicula. Donec ornare orci vitae nibh porta commodo. Suspendisse et lacus ut mauris rhoncus euismod. Maecenas ex justo, semper ac porttitor id, auctor nec dui. Nam tempus, lorem et viverra rutrum, felis arcu finibus mi, in porta massa mi eu lacus. Sed faucibus nulla quis nunc ultrices, quis cursus mi aliquet.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(3, 'Cras mollis mi sit amet libero elementum, ac mollis dui porta', 'Nulla sagittis enim velit, quis gravida ligula varius vitae. Curabitur neque lorem, lobortis ac congue non, volutpat vel purus. Mauris sed tortor quam. Suspendisse in erat in elit sollicitudin molestie molestie vehicula ipsum. Curabitur at eros at massa elementum tincidunt. Quisque cursus aliquet ex consectetur bibendum. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus ac nulla mattis leo fringilla tempor. Nam aliquam eu ligula vitae ultricies. Mauris bibendum urna rutrum porta egestas.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(3, 'Sed faucibus purus et condimentum lobortis', 'Maecenas sit amet ligula lacus. Phasellus imperdiet dignissim dolor, vitae tincidunt augue cursus ac. Vivamus facilisis lacus a elit vestibulum tristique. Phasellus suscipit, turpis at suscipit viverra, ante lorem egestas turpis, non tincidunt ante purus at lorem. Integer molestie varius ligula, vel accumsan lectus vestibulum sit amet. Nam sed urna quis elit tempor eleifend auctor eu odio. Donec ullamcorper consectetur convallis. Aenean lacinia purus ac metus dapibus, sed molestie ligula posuere. Curabitur bibendum est at eros hendrerit convallis. Quisque elementum eget felis id ornare. Quisque blandit mollis faucibus.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(3, 'Curabitur at ante ultricies, sagittis dolor a, dictum est', ' Nulla ullamcorper mi nec dui pretium suscipit. Donec et tortor placerat, scelerisque mi a, interdum diam. Nam vulputate egestas lacinia. Fusce sit amet dui porta, rutrum metus et, tincidunt ligula. Sed vitae leo non ex pellentesque dignissim in ut eros. Vivamus iaculis neque odio, non dictum purus mollis sit amet. In hendrerit tempor metus, nec dignissim magna condimentum non. Praesent at tortor id turpis efficitur imperdiet id eget lorem. Quisque sollicitudin et nisl in molestie.

Nam condimentum felis sed sapien euismod, non iaculis metus sollicitudin. Donec sed tincidunt augue, eu consequat est. Proin non magna dignissim, fringilla magna ut, accumsan enim. Pellentesque tincidunt, nisl at ornare porttitor, massa risus tincidunt sem, id sagittis sapien nibh vitae augue. Cras in placerat metus. Praesent nulla purus, accumsan vel dignissim quis, tincidunt et diam. Phasellus eu dolor et eros aliquet molestie eu in libero. Nam vel tempor neque, quis egestas elit. Ut dapibus tellus et nulla vehicula, et ultrices orci eleifend. Ut malesuada mauris at sodales vulputate. Fusce porttitor, lorem et volutpat tristique, lectus augue pharetra nibh, a ornare massa arcu quis sapien. In eu tortor at sem sodales hendrerit id vitae nisl. Sed quis imperdiet augue, eget eleifend odio. Phasellus laoreet in est et cursus.', CURRENT_TIMESTAMP);

INSERT INTO post (user_id, title, body, created) VALUES(3, 'Fusce at tellus nec sapien pretium tincidunt', ' Ut placerat rutrum ornare. Vestibulum dapibus nibh dolor, consectetur tempus dolor mollis vitae. Donec in rhoncus tellus. Curabitur eget lorem ante. Sed eget nisi ultricies, sodales risus a, sodales felis. Nulla faucibus nunc vel sem facilisis, et egestas turpis lobortis. Proin eros dui, volutpat non vehicula quis, pharetra at lacus. Sed fringilla, diam at mollis facilisis, urna lacus auctor nunc, id efficitur dui dolor in erat. Ut molestie pulvinar posuere. Aliquam congue mauris at turpis eleifend vestibulum. Duis orci nulla, fermentum nec ante ac, rutrum ornare leo. Morbi a nibh dolor. Nulla vel sollicitudin nisi, sed volutpat nisi. In dictum ac nisi at ultrices. Morbi massa enim, lobortis nec lacus non, blandit luctus neque.

Phasellus vel dapibus libero. Quisque in ante leo. Sed in tellus et lorem viverra auctor. Proin vitae enim mi. Duis viverra sollicitudin lectus ac porttitor. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi fringilla erat egestas, scelerisque lorem sed, congue turpis. Maecenas faucibus varius lorem, quis interdum enim faucibus eu. Sed sed quam pellentesque, lobortis est et, laoreet urna. Aenean in tempus neque, sit amet porttitor nulla. Phasellus vulputate vehicula lorem, ut bibendum ante posuere eget. Duis volutpat dolor nec varius mollis. Proin sit amet consectetur neque, at rhoncus dui. Donec consectetur ornare elementum. Ut vestibulum mi orci, in tincidunt velit cursus sed.', CURRENT_TIMESTAMP);

-- post-votes data
INSERT INTO post_vote (post_id, user_id, vote) VALUES(1, 1, 1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(1, 2, 1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(2, 1, 1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(2, 2, 1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(3, 1, -1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(3, 2, -1);
INSERT INTO post_vote (post_id, user_id, vote) VALUES(4, 1, 1);

-- post-categories data
INSERT INTO post_category (post_id, category_id) VALUES(1, 1);
INSERT INTO post_category (post_id, category_id) VALUES(1, 2);
INSERT INTO post_category (post_id, category_id) VALUES(1, 3);
INSERT INTO post_category (post_id, category_id) VALUES(1, 4);
INSERT INTO post_category (post_id, category_id) VALUES(1, 5);
INSERT INTO post_category (post_id, category_id) VALUES(2, 1);
INSERT INTO post_category (post_id, category_id) VALUES(2, 4);
INSERT INTO post_category (post_id, category_id) VALUES(2, 5);
INSERT INTO post_category (post_id, category_id) VALUES(3, 2);
INSERT INTO post_category (post_id, category_id) VALUES(4, 3);
INSERT INTO post_category (post_id, category_id) VALUES(5, 5);
INSERT INTO post_category (post_id, category_id) VALUES(6, 5);
INSERT INTO post_category (post_id, category_id) VALUES(6, 3);
INSERT INTO post_category (post_id, category_id) VALUES(6, 4);
INSERT INTO post_category (post_id, category_id) VALUES(7, 1);
INSERT INTO post_category (post_id, category_id) VALUES(8, 2);
INSERT INTO post_category (post_id, category_id) VALUES(8, 3);
INSERT INTO post_category (post_id, category_id) VALUES(9, 4);
INSERT INTO post_category (post_id, category_id) VALUES(9, 5);
INSERT INTO post_category (post_id, category_id) VALUES(9, 1);
INSERT INTO post_category (post_id, category_id) VALUES(9, 2);

-- comments data
INSERT INTO comment (body, user_id, post_id, created) VALUES('Aenean et orci tempor, ullamcorper magna et, cursus magna', 1, 1, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Nullam vel est a risus efficitur cursus', 1, 1, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Quisque elementum arcu eu est euismod bibendum', 2, 2, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Pellentesque quis tortor pretium, commodo est eget, feugiat nisl', 2, 3, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Fusce viverra est sed augue tincidunt, et dictum diam cursus', 1, 2, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Pellentesque in est vel est vestibulum ornare', 1, 4, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Fusce lobortis risus placerat, ornare orci ac, blandit magna', 1, 3, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Mauris scelerisque erat et tortor euismod, vitae aliquet lorem pulvinar', 3, 4, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Nullam at lorem at eros rhoncus lacinia a eget dolor', 3, 5, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Ut sit amet dolor sit amet orci auctor tempor', 3, 6, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Aenean at libero sit amet urna dapibus congue', 2, 7, CURRENT_TIMESTAMP);
INSERT INTO comment (body, user_id, post_id, created) VALUES('Donec sollicitudin neque auctor ex euismod accumsan', 1, 8, CURRENT_TIMESTAMP);

-- Create a table to store comment-votes
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(1, 1, 1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(1, 2, 1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(2, 1, 1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(2, 2, -1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(3, 2, 1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(3, 1, -1);
INSERT INTO comment_vote (comment_id, user_id, vote) VALUES(4, 1, 1);
