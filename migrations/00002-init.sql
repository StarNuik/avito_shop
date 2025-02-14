create table Users (
    Id bigint primary key generated always as identity,
    Username text not null,
    PasswordHash text not null,
    Coins bigint not null
);
create index idx_Users_Username on Users(Username);

create table Transfers (
    Id bigint primary key generated always as identity,
    Delta bigint not null,
    FromUser bigint references Users(Id),
    ToUser bigint references Users(Id)
);
create index idx_Transfers_FromUser on Transfers(FromUser);
create index idx_Transfers_ToUser on Transfers(FromUser);

create table Inventory (
    Id bigint primary key generated always as identity,
    Name text not null,
    Price bigint not null
);
create index idx_Inventory_Name on Inventory(Name);

create table Purchases (
    Id bigint primary key generated always as identity,
    Price bigint not null,
    Item bigint references Inventory(Id),
    UserId bigint references Users(Id)
);
create index idx_Purchases_UserId on Purchases(UserId);
