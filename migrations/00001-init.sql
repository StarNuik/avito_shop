create table Users (
    Id bigint primary key,
    Username text not null,
    PasswordHash text not null
);
-- TODO: index Users(Username)

create table BalanceOperations (
    Id bigint primary key,
    User bigint references Users(Id),
    Delta bigint not null,
    Result bigint not null
);

create table Transfers (
    Id bigint primary key,
    SourceOp bigint references BalanceOperations(Id),
    TargetOp bigint references BalanceOperations(Id)
);

create table Inventory (
    Id bigint primary key,
    Name text not null,
    Price bigint not null
);
-- TODO: index Inventory(Name)

create table Purchases (
    Id bigint primary key,
    Item bigint references Inventory(Id),
    User bigint references Users(Id),
    Operation bigint references BalanceOperations(Id)
);

