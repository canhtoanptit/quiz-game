package com.toannc.quiz.usermngservice.dtos;

import lombok.Getter;
import lombok.Setter;

/**
 * @author toannguyen
 * @since 04/10/2024
 */
@Getter
@Setter
public class AuthRequest {

    private String username;
    private String password;
}
