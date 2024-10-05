package com.toannc.quiz.usermngservice.repositories;

import com.toannc.quiz.usermngservice.entities.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

/**
 * @author toannguyen
 * @since 04/10/2024
 */
@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    Optional<User> findByUsername(String username);
}
