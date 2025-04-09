


create table  UserOnPlatform (
    id uuid primary key,
    full_name varchar(100) not null ,
    email varchar(100) not null,
    password varchar(88) not null,
    verified bool  not null ,
    createdAt timestamp not null,
    walletCreated bool not null ,
    walletCreatedTime timestamp not null
);

create table  Skills (
    id uuid primary key ,
    skillTag varchar(100) not null
);

create  table User_Skill (
    skill_id uuid  not null ,
    user_id uuid not null,
    primary key (user_id, skill_id),
    foreign key (user_id) references UserOnPlatform(id),
    foreign key (skill_id) references Skills(id)
);


create  table UserBiometrics (
    id uuid primary key  not null ,
    user_id uuid not null ,
    biometricsHash varchar(200) not null ,
    foreign key (user_id) references UserOnPlatform(id)
);


create table  Services (
    id uuid not null ,
    service_name varchar(100) not null ,
    description text not null ,
    offeredBy uuid not null ,
    createdAt timestamp default  current_time,
    updatedAt timestamp,
    available bool  not null,
    minimum  float8 not null ,
    maximum float8 not null ,
    foreign key  (offeredBy) references UserOnPlatform(id)
);

create table Category (
    id uuid not null ,
    category_name varchar(30) not null
);

create  table  Service_Category (
    category_id uuid not null ,
    service_id uuid not null ,
    primary key (category_id, service_id),
    foreign key (category_id) references Category(id),
    foreign key  (service_id) references  Services(id)
);
create table service_image  (
    url varchar(200) not null,
    description text not null,
    ownedBy uuid not null ,
    belongedTo uuid not null,
    primary key (ownedBy, belongedTo),
    foreign key (ownedBy) references UserOnPlatform(id),
    foreign key  (belongedTo) references  Services(id)
);


create  table orders (
    id uuid  not null ,
    buyer uuid not null ,
    seller uuid not null,
    underlying_service_id uuid not null ,
    price float8 not null,
    createdAt timestamp not null ,
    completedAt timestamp not null ,
    status varchar(100) not null,
    payment_method varchar(100) not null,
    transaction_id varchar(255),
    primary key ( buyer, seller, underlying_service_id),
    foreign key (buyer) references UserOnPlatform(id),
    foreign key (seller) references  UserOnPlatform(id),
    foreign key (underlying_service_id) references Services(id)
);
CREATE TABLE Payments (
                          id UUID PRIMARY KEY,
                          order_id UUID NOT NULL,
                          payment_method VARCHAR(100) NOT NULL,
                          transaction_id VARCHAR(255),
                          FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE Chat (
                      chat_room_id UUID NOT NULL,
                      message_id UUID PRIMARY KEY,
                      user_id UUID NOT NULL,
                      sender_id UUID NOT NULL,
                      message TEXT NOT NULL,
                      message_type VARCHAR(50) DEFAULT 'text',
                      timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                      FOREIGN KEY (user_id) REFERENCES UserOnPlatform(id) ON DELETE CASCADE,
                      FOREIGN KEY (sender_id) REFERENCES UserOnPlatform(id) ON DELETE CASCADE
);

