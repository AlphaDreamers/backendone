create table "Orders"
(
    id             uuid default gen_random_uuid() not null
        primary key,
    order_number   text                           not null,
    price          numeric                        not null,
    payment_method text                           not null,
    status         text default 'PENDING'::text   not null,
    transaction_id text,
    requirements   text,
    package_id     uuid                           not null,
    seller_id      uuid                           not null,
    buyer_id       uuid                           not null,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    completed_at   timestamp with time zone,
    due_date       timestamp with time zone
);

alter table "Orders"
    owner to postgres;

create unique index "idx_Orders_order_number"
    on "Orders" (order_number);

create table "Badge"
(
    id          uuid default gen_random_uuid() not null
        primary key,
    label       text                           not null,
    icon        text                           not null,
    color       text                           not null,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "Badge"
    owner to postgres;

create table "User"
(
    id                  uuid    default gen_random_uuid() not null
        primary key,
    "firstName"         text                              not null,
    "lastName"          text                              not null,
    email               text                              not null,
    password            text                              not null,
    verified            boolean default false             not null,
    "twoFactorVerified" boolean default false             not null,
    username            text                              not null,
    avatar              text,
    country             text                              not null,
    "walletCreated"     boolean default false             not null,
    "walletCreatedTime" timestamp with time zone,
    "createdAt"         timestamp with time zone,
    "updatedAt"         timestamp with time zone
);

alter table "User"
    owner to postgres;

create unique index "idx_User_email"
    on "User" (email);

create table "UserBadge"
(
    id           uuid    default gen_random_uuid() not null
        primary key,
    "userId"     uuid                              not null
        constraint "fk_User_user_badges"
            references "User",
    "badgeId"    uuid                              not null
        constraint "fk_UserBadge_badge"
            references "Badge"
            on update cascade on delete cascade,
    tier         text    default 'BRONZE'::text    not null,
    "isFeatured" boolean default false             not null,
    "createdAt"  timestamp with time zone,
    "updatedAt"  timestamp with time zone
);

alter table "UserBadge"
    owner to postgres;

create table "Skill"
(
    id          uuid default gen_random_uuid() not null
        primary key,
    label       text                           not null,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "Skill"
    owner to postgres;

create unique index "idx_Skill_label"
    on "Skill" (label);

create table user_skills
(
    "skillId"   uuid                  not null
        constraint "fk_Skill_users"
            references "Skill",
    "userId"    uuid                  not null
        constraint "fk_User_skills"
            references "User",
    level       bigint  default 1     not null,
    endorsed    boolean default false not null,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    primary key ("skillId", "userId")
);

alter table user_skills
    owner to postgres;

create table "Biometrics"
(
    id           uuid    default gen_random_uuid() not null
        primary key,
    type         text                              not null,
    value        text                              not null,
    "isVerified" boolean default false             not null,
    "userId"     uuid                              not null
        constraint "fk_User_biometrics"
            references "User",
    "createdAt"  timestamp with time zone,
    "updatedAt"  timestamp with time zone
);

alter table "Biometrics"
    owner to postgres;

create table "GigTag"
(
    id          uuid default gen_random_uuid() not null
        primary key,
    label       text                           not null,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "GigTag"
    owner to postgres;

create unique index "idx_GigTag_label"
    on "GigTag" (label);

create table "Category"
(
    id          uuid    default gen_random_uuid() not null
        primary key,
    label       text                              not null,
    slug        text                              not null,
    "isActive"  boolean default true              not null,
    "sortOrder" bigint  default 0                 not null,
    "parentId"  uuid
        constraint "fk_Category_children"
            references "Category",
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "Category"
    owner to postgres;

create unique index "idx_Category_slug"
    on "Category" (slug);

create unique index "idx_Category_label"
    on "Category" (label);

create table "Gig"
(
    id              uuid    default gen_random_uuid() not null
        primary key,
    title           text                              not null,
    description     text                              not null,
    "isActive"      boolean default true              not null,
    "viewCount"     bigint  default 0                 not null,
    "averageRating" numeric default 0                 not null,
    "ratingCount"   bigint  default 0                 not null,
    "categoryId"    uuid                              not null
        constraint "fk_Category_gigs"
            references "Category",
    "sellerId"      uuid                              not null
        constraint "fk_User_gigs"
            references "User",
    "createdAt"     timestamp with time zone,
    "updatedAt"     timestamp with time zone
);

alter table "Gig"
    owner to postgres;

create table _gig_to_gig_tags
(
    a uuid default gen_random_uuid() not null
        constraint fk__gig_to_gig_tags_gig
            references "Gig",
    b uuid default gen_random_uuid() not null
        constraint fk__gig_to_gig_tags_gig_tag
            references "GigTag",
    primary key (a, b)
);

alter table _gig_to_gig_tags
    owner to postgres;

create table "RegistrationToken"
(
    code        text                     not null
        primary key,
    email       text                     not null,
    "expiresAt" timestamp with time zone not null,
    "userId"    uuid
        constraint "fk_RegistrationToken_user"
            references "User"
            on update cascade on delete set null,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "RegistrationToken"
    owner to postgres;

create unique index "idx_RegistrationToken_email"
    on "RegistrationToken" (email);

create table "GigImage"
(
    id          uuid    default gen_random_uuid() not null
        primary key,
    url         text                              not null,
    alt         text,
    "isPrimary" boolean default false             not null,
    "sortOrder" bigint  default 0                 not null,
    "gigId"     uuid                              not null
        constraint "fk_Gig_images"
            references "Gig",
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "GigImage"
    owner to postgres;

create table "GigPackage"
(
    id             uuid    default gen_random_uuid() not null
        primary key,
    title          text                              not null,
    description    text                              not null,
    price          numeric                           not null,
    "deliveryTime" bigint                            not null,
    revisions      bigint  default 1                 not null,
    "isActive"     boolean default true              not null,
    "gigId"        uuid                              not null
        constraint "fk_Gig_packages"
            references "Gig",
    "createdAt"    timestamp with time zone,
    "updatedAt"    timestamp with time zone
);

alter table "GigPackage"
    owner to postgres;

create table "GigPackageFeature"
(
    id             uuid    default gen_random_uuid() not null
        primary key,
    title          text                              not null,
    description    text,
    included       boolean default true              not null,
    "gigPackageId" uuid                              not null
        constraint "fk_GigPackage_features"
            references "GigPackage",
    "createdAt"    timestamp with time zone,
    "updatedAt"    timestamp with time zone
);

alter table "GigPackageFeature"
    owner to postgres;

create table "Order"
(
    id             uuid default gen_random_uuid() not null
        primary key,
    order_number   text                           not null,
    price          numeric                        not null,
    payment_method text                           not null,
    status         text default 'PENDING'::text   not null,
    transaction_id text,
    requirements   text,
    package_id     uuid                           not null,
    seller_id      uuid                           not null
        constraint "fk_User_orders_as_seller"
            references "User",
    buyer_id       uuid                           not null
        constraint "fk_User_orders_as_buyer"
            references "User",
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    completed_at   timestamp with time zone,
    due_date       timestamp with time zone
);

alter table "Order"
    owner to postgres;

create unique index "idx_Order_order_number"
    on "Order" (order_number);

create table "Review"
(
    id          uuid    default gen_random_uuid() not null
        primary key,
    title       text                              not null,
    content     text,
    rating      bigint  default 5                 not null,
    "isPublic"  boolean default true              not null,
    "authorId"  uuid                              not null
        constraint "fk_Review_author"
            references "User"
            on update cascade on delete restrict,
    "orderId"   uuid                              not null
        constraint "fk_Review_order"
            references "Order"
            on update cascade on delete restrict,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone
);

alter table "Review"
    owner to postgres;

create unique index "idx_Review_order_id"
    on "Review" ("orderId");

create index "idx_Review_author_id"
    on "Review" ("authorId");

