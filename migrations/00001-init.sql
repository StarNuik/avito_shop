create table Users (
    Id bigint primary key,
    Username text not null,
    PasswordHash text not null
);
-- TODO: index Users on Username

create table Transactions (
    Id bigint primary key,
    Delta bigint not null,
    Result bigint not null,
    UserFrom bigint references Users(Id),
    UserTo bigint references Users(Id)
);

create table Inventory (
    Id bigint primary key,
    Name text not null,
    Price bigint not null
);

create table Purchases (
    Id bigint primary key,
    Item bigint references Inventory(Id),
    User bigint references Users(Id),
    Transaction bigint references Transactions(Id)
);

