    CREATE TABLE UserProfile(
        profileID      INTEGER PRIMARY KEY,
        firstName   VARCHAR(50)      NOT NULL,
        lastName    VARCHAR(50)     NOT NULL,
        accountID  INTEGER  NOT NULL

    ); 

    ALTER TABLE UserProfile
    ADD CONSTRAINT FK_account
    FOREIGN KEY (accountID) REFERENCES Account(accountID);

INSERT INTO userProfile (firstName, lastName, email, password) 
VALUES ('Elijah', 'Rose', 'elijahr0001@gmail.com', 'test123'); 


    CREATE TABLE CourierProfile(
        profileID   INTEGER PRIMARY KEY,
        firstName   VARCHAR(50)      NOT NULL,
        lastName    VARCHAR(50)     NOT NULL,
        overallRating INTEGER
        accountID     INTEGER  NOT NULL
    ); 

ALTER TABLE CourierProfile
ADD CONSTRAINT FK_account
FOREIGN KEY (accountID) REFERENCES Account(accountID);
    

INSERT INTO courieraccount (firstName, lastName, email, password) 
VALUES ('Elijah', 'Rose', 'poisondart14@gmail.com', 'test123');  




INSERT INTO Delivery (userID, deliveryType, collectionAddress, deliveryAddress, expiryDate, accepted, status) 
VALUES (1,'Food', '18 Victora Lane, Folkestone, kent, FT67 889', '25 Pleasant view, Folkestone, kent, FTP7 8U9', NOW(), false, 'Waiting');  


    CREATE TABLE Delivery(
        deliveryID  INTEGER     PRIMARY KEY,
        userID      INTEGER     REFERENCES      UserProfile(profileID),
        deliveryType          VARCHAR(50)         NOT NULL,
        collectionAddress     VARCHAR(255)        NOT NULL,
        deliveryAddress       VARCHAR(255)        NOT NULL,
        expiryDate            TIMESTAMP           NOT NULL, 
        extraDetails          VARCHAR(255),
        accepted              BOOLEAN             NOT NULL, 
        status                VARCHAR(50)         NOT NULL,
        collectionLat         FLOAT               NOT NULL, 
        collectionLon         FLOAT               NOT NULL
        deliveryLat           FLOAT               NOT NULL,
        deliveryLon           FLOAT               NOT NULL
    ); 

    ALTER TABLE Delivery ALTER COLUMN collectionAddress  TYPE varchar (255);


    CREATE TABLE Accepted_delivery(
        deliveryID    INTEGER     REFERENCES      Delivery(deliveryID),
        profileID     INTEGER     REFERENCES      CourierProfile(profileID),
        deliveryDate  VARCHAR(50)                NOT NULL
    ); 


    CREATE TABLE Review(
        reviewID   SERIAL     PRIMARY KEY,
        userID     SERIAL     REFERENCES      UserAccount(userID),
        courierID  SERIAL     REFERENCES      courierAccount(courierID),
        title                 VARCHAR(255)    NOT NULL,
        content               VARCHAR(300)    NOT NULL,
        rating                INTEGER         NOT NULL,
        date                  TIMESTAMP       NOT NULL, 
        userVisible           BOOLEAN         NOT NULL
    ); 


 CREATE TABLE Account(
        accountID   INTEGER PRIMARY KEY,
        email       VARCHAR(255)    NOT NULL UNIQUE,
        password    VARCHAR(255)    NOT NULL
); 




