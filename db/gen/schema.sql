create table
    users (
        _id bigint primary key,
        email varchar(280) not null,
        user_status int not null,
        user_name varchar(50) not null,
        password char(60) not null,
        phone varchar(16) not null,
        bio varchar(256) not null,
        xp bigint not null,
        level int not null,
        tier int not null,
        user_rank int not null,
        coffee bigint not null,
        otp varchar(64),
        otp_validated boolean,
        stripe_user varchar(280),
        stripe_subscription varchar(280),
        first_name varchar(280),
        last_name varchar(280),
        auth_role int not null,
        gitea_id bigint not null,
        external_auth char(21),
        created_at datetime not null,
        workspace_settings json not null,
        follower_count bigint not null,
        start_user_info json not null,
        highest_score bigint not null,
        timezone varchar(50) not null,
        avatar_settings json not null,
        broadcast_threshold bigint not null,
        avatar_reward bigint,
        encrypted_service_key varbinary(500) not null,
        stripe_account varchar(280),
        exclusive_agreement boolean not null default false,
        reset_token varchar(500),
        hasBroadcast boolean not null default false,
        holiday_themes boolean not null default true,
        is_ephemeral boolean not null default false
    );

create table
    attempt (
        _id bigint primary key,
        post_title varchar(150) not null,
        description varchar(500) not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        updated_at datetime not null,
        repo_id bigint not null,
        author_tier int not null,
        coffee bigint not null,
        post_id bigint not null,
        closed boolean not null,
        success boolean not null,
        closed_date datetime,
        tier int not null,
        parent_attempt bigint,
        workspace_settings json,
        post_type int not null
    );

create table
    attempt_awards (
        attempt_id bigint not null,
        award_id bigint not null
    );

create table
    award (
        _id bigint not null primary key,
        types int not null,
        award varchar(280) not null
    );

create table
    broadcast_event (
        _id bigint primary key,
        user_id bigint not null,
        user_name varchar(50) not null,
        message longtext not null,
        broadcast_type int not null,
        time_posted datetime not null
    );

create table
    coffee (
        _id bigint primary key,
        created_at datetime not null,
        updated_at datetime not null,
        post_id bigint,
        attempt_id bigint,
        user_id bigint not null,
        discussion_id bigint not null
    );

create table
    comment (
        _id bigint,
        body longtext not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        author_tier int not null,
        coffee bigint not null,
        discussion_id bigint not null,
        leads boolean not null default false,
        revision int not null,
        discussion_level int not null,
        primary key (_id, revision)
    );

create table
    comment_awards (
        comment_id bigint not null,
        award_id bigint not null,
        revision int not null,
        primary key (comment_id, award_id, revision)
    );

create table
    discussion (
        _id bigint,
        body longtext not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        updated_at datetime not null,
        author_tier int not null,
        coffee bigint not null,
        post_id bigint not null,
        title varchar(280) not null,
        leads boolean not null default false,
        revision int not null,
        discussion_level int not null,
        primary key (_id, revision)
    );

create table
    discussion_awards (
        discussion_id bigint not null,
        award_id bigint not null,
        revision int not null,
        primary key (
            discussion_id,
            award_id,
            revision
        )
    );

create table
    discussion_tags (
        discussion_id bigint not null,
        tag_id bigint not null,
        revision int not null,
        primary key (
            discussion_id,
            tag_id,
            revision
        )
    );

create table
    discussion_up_vote (
        discussion_id bigint not null,
        up_vote_id bigint not null,
        user_id bigint not null,
        primary key (
            discussion_id,
            up_vote_id,
            user_id
        )
    );

create table
    follower (
        follower bigint not null,
        following bigint not null,
        primary key (follower, following)
    );

create table
    friend_requests (
        _id bigint primary key not null,
        user_id bigint not null,
        user_name varchar(50) not null,
        friend bigint not null,
        friend_name varchar(50) not null,
        response boolean null,
        date datetime not null,
        notification_id bigint not null
    );

create table
    friends (
        _id bigint primary key not null,
        user_id bigint not null,
        user_name varchar(50) not null,
        friend bigint not null,
        friend_name varchar(50) not null,
        date datetime not null
    );

create table
    implicit_rec (
        _id bigint primary key not null,
        user_id bigint not null,
        post_id bigint not null,
        session_id binary(16) not null,
        implicit_action int not null,
        created_at timestamp not null,
        user_tier_at_action int not null
    );

create table
    nemesis (
        _id bigint primary key not null,
        antagonist_id bigint not null,
        antagonist_name varchar(50) not null,
        antagonist_towers_captured bigint not null,
        protagonist_id bigint not null,
        protagonist_name varchar(50) not null,
        protagonist_towers_captured bigint not null,
        time_of_villainy datetime not null,
        victor bigint,
        is_accepted boolean not null,
        end_time datetime
    );

create table
    nemesis_history (
        _id bigint primary key not null,
        match_id bigint not null,
        antagonist_id bigint not null,
        protagonist_id bigint not null,
        protagonist_towers_held bigint not null,
        antagonist_towers_held bigint not null,
        protagonist_total_xp bigint not null,
        antagonist_total_xp bigint not null,
        is_alerted boolean not null,
        created_at datetime not null
    );

create table
    notification (
        _id bigint primary key not null,
        user_id bigint not null,
        message longtext not null,
        notification_type int not null,
        created_at datetime not null,
        acknowledged boolean not null default false,
        interacting_user_id bigint
    );

create table
    post (
        _id bigint primary key not null,
        title varchar(150) not null,
        description varchar(500) not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        updated_at datetime not null,
        repo_id bigint not null,
        top_reply bigint,
        tier int not null,
        coffee bigint not null,
        post_type int not null,
        views bigint not null,
        completions bigint not null,
        attempts bigint not null,
        published boolean not null,
        visibility int not null,
        challenge_cost varchar(16),
        stripe_price_id varchar(50),
        workspace_config bigint not null,
        workspace_config_revision int not null,
        workspace_settings json,
        leads boolean not null default false,
        embedded boolean not null default false,
        deleted boolean not null default false,
        exclusive_description varchar(500),
        share_hash varchar(16),
        estimated_tutorial_time bigint
    );

create table
    post_awards (
        post_id bigint not null,
        award_id bigint not null,
        primary key (post_id, award_id)
    );

create table
    post_tags (
        post_id bigint not null,
        tag_id bigint not null,
        primary key (post_id, tag_id)
    );

create table
    post_langs (
        post_id int not null,
        lang_id bigint not null,
        primary key (post_id, lang_id)
    );

create table
    recommended_post (
        _id bigint primary key not null,
        user_id bigint not null,
        post_id bigint not null,
        type int not null,
        reference_id bigint not null,
        score float not null,
        created_at timestamp not null,
        expires_at timestamp not null,
        reference_tier int not null,
        accepted boolean not null default false,
        views bigint not null default 0
    );

create table
    rewards (
        _id bigint primary key not null,
        color_palette varchar(280) not null,
        render_in_front boolean not null,
        name varchar(280) not null
    );

create table
    user_rewards_inventory (
        reward_id bigint not null,
        user_id bigint not null,
        primary key (reward_id, user_id)
    );

create table
    search_rec (
        _id bigint primary key not null,
        user_id bigint not null,
        query text not null,
        selected_post_id bigint,
        selected_post_name text,
        created_at timestamp not null
    );

create table
    search_rec_posts (
        search_id bigint not null,
        post_id bigint not null,
        primary key (search_id, post_id)
    );

create table
    user_stats (
        _id bigint primary key not null,
        user_id bigint not null,
        challenges_completed int,
        streak_active boolean,
        current_streak int,
        longest_streak int not null,
        total_time_spent bigint not null,
        avg_time bigint not null,
        days_on_platform int not null,
        days_on_fire int not null,
        streak_freezes int not null,
        streak_freeze_used boolean not null,
        xp_gained bigint not null,
        date datetime not null,
        expiration datetime not null,
        closed boolean not null default false
    );

create table
    user_daily_usage (
        user_id bigint not null,
        start_time datetime not null,
        end_time datetime,
        open_session int not null,
        date datetime not null,
        primary key (user_id, start_time, date)
    );

create table
    stats_xp (
        stats_id bigint primary key not null,
        expiration datetime not null
    );

create table
    tag (
        _id bigint primary key not null,
        value varchar(50) not null,
        official boolean not null,
        usage_count bigint not null
    );

create table
    thread_comment (
        _id bigint primary key not null,
        body longtext not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        author_tier int not null,
        coffee bigint not null,
        comment_id bigint not null,
        leads boolean not null default false,
        revision int not null,
        discussion_level int not null
    );

create table
    thread_reply (
        _id bigint not null,
        body longtext not null,
        author varchar(16) not null,
        author_id bigint not null,
        created_at datetime not null,
        author_tier int not null,
        coffee bigint not null,
        thread_comment_id bigint not null,
        revision int not null,
        discussion_level int not null,
        primary key (_id, revision)
    );

create table
    up_vote (
        _id bigint primary key not null,
        discussion_type int not null,
        discussion_id bigint not null,
        user_id bigint not null
    );

create table
    user_active_times (
        _id bigint primary key not null,
        user_id bigint not null,
        start_time datetime not null,
        end_time datetime not null
    );

create table
    user_badges (
        user_id bigint not null,
        badge_id bigint not null
    );

create table
    user_saved_posts (
        user_id bigint not null,
        post_id bigint not null
    );

create table
    user_free_premium (
        _id bigint primary key not null,
        user_id bigint not null,
        start_date datetime not null,
        end_date datetime not null,
        length varchar(50) not null
    );

create table
    user_session_key (
        _id bigint primary key not null,
        _key binary(128) not null,
        expiration datetime not null
    );

create table
    workspaces (
        _id bigint primary key not null,
        code_source_id bigint not null,
        code_source_type int not null,
        repo_id bigint not null,
        created_at datetime not null,
        owner_id bigint not null,
        template_id bigint not null,
        expiration datetime not null,
        commit
            varchar(64) not null,
            state bigint not null,
            init_state int not null,
            init_failure json,
            last_state_update datetime not null,
            workspace_settings json not null,
            over_allocated json,
            ports json,
            is_ephemeral boolean not null default false
    );

create table
    workspace_agent (
        _id bigint not null primary key,
        created_at datetime not null,
        updated_at datetime not null,
        first_connect datetime,
        last_connect datetime,
        last_disconnect datetime,
        last_connected_node bigint,
        disconnect_count int not null,
        state int not null,
        workspace_id bigint not null,
        version varchar(30) not null,
        owner_id bigint not null,
        secret binary(16) not null
    );

create table
    workspace_agent_stats (
        _id bigint not null primary key,
        agent_id bigint not null,
        workspace_id bigint not null,
        timestamp datetime not null,
        conns_by_proto json not null,
        num_comms bigint not null,
        rx_packets bigint not null,
        rx_bytes bigint not null,
        tx_packets bigint not null,
        tx_bytes bigint not null
    );

create table
    workspace_config (
        _id bigint not null,
        title varchar(50) not null,
        description varchar(500) not null,
        content longtext not null,
        author_id bigint not null,
        revision int not null,
        official boolean not null,
        primary key (_id, revision)
    );

create table
    workspace_config_tags (
        cfg_id bigint not null,
        tag_id bigint not null,
        revision int not null,
        primary key (cfg_id, tag_id, revision)
    );

create table
    workspace_config_langs (
        cfg_id int not null,
        lang_id bigint not null,
        revision int not null,
        primary key (cfg_id, lang_id, revision)
    );

create table
    if not exists xp_boosts (
        _id bigint not null primary key,
        user_id bigint not null,
        end_date datetime
    );

create table
    if not exists zookies (
        r_id varchar(36) not null,
        s_id varchar(36) not null,
        z varchar(64) not null,
        time datetime not null,
        end_date datetime
    );

create table
    if not exists exclusive_content_purchases (
        user_id bigint not null,
        post bigint not null,
        date datetime not null,
        primary key (user_id, post)
    );

create table
    if not exists report_issue (
        _id bigint not null primary key,
        user_id bigint not null,
        date datetime not null,
        issue longtext not null,
        page varchar(36) not null
    );

create table
    if not exists xp_reasons (
        _id bigint not null primary key,
        user_id bigint not null,
        date timestamp not null,
        reason varchar(280) not null,
        xp bigint not null
    );

create table
    if not exists chat (
        _id bigint not null primary key,
        name varchar(280) not null,
        type int not null,
        last_message_time datetime
    );

create table
    if not exists chat_users (
        chat_id bigint not null,
        user_id bigint not null,
        created_at datetime not null,
        primary key (chat_id, user_id)
    );

create table
    if not exists chat_messages (
        _id bigint not null primary key,
        chat_id bigint not null,
        author_id bigint not null,
        author varchar(280) not null,
        message varchar(500) not null,
        created_at datetime(6) not null,
        revision bigint not null,
        type int not null
    );

create table
    if not exists curated_post (
        _id bigint not null primary key,
        post_id bigint not null,
        post_language int not null
    );

create table
    if not exists curated_post_type (
        curated_id bigint not null,
        proficiency_type int not null,
        primary key (curated_id, type)
    );

create table
    if not exists ephemeral_shared_workspaces (
        workspace_id BIGINT NOT NULL,
        ip BIGINT NOT NULL,
        date DATETIME NOT NULL,
        primary key (workpace_id, ip)
    );

create table
    if not exists web_tracking (
    `_id` bigint not null primary key,
    `user_id` bigint default NULL,
    `ip` bigint not null,
    `host` varchar(255) not null,
    `event` varchar(255) not null,
    `timestamp` datetime not null,
    `timespent` bigint default NULL,
    `path` varchar(255) not null,
    `lattitude` double,
    `longitude` double,
    `metadata` json
)