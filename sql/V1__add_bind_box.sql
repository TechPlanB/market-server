alter table nft_tokens
    add blind_box bool default false null comment '是否是盲盒 0：否 1：是' after owner;

alter table nft_listeners
    add blind_box bool default false null comment '是否是盲盒 0：否 1：是' after symbol;