package com.sidekick.dataservice.controller;

import com.sidekick.dataservice.component.EncryptUtils;
import com.sidekick.dataservice.component.ResultMap;
import com.sidekick.dataservice.enumeration.ErrCode;
import com.sidekick.dataservice.enumeration.Conf;
import com.sidekick.dataservice.mapper.UserMapper;
import com.sidekick.dataservice.model.User;
import com.sidekick.dataservice.view.request.UserPassword;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import java.security.NoSuchAlgorithmException;
import java.util.List;

@RequestMapping("/v1.0/user")
@RestController
public class UserController {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    @Autowired
    private UserMapper userMapper;

    @Autowired
    private EncryptUtils encryptUtils;

    @Autowired
    private ResultMap resultMap;

    @RequestMapping("/add")
    public Object addUser(HttpServletRequest request, @RequestBody User user) throws NoSuchAlgorithmException {
        logger.info("/v1.0/user/add");

        if (user.getName() == null ||
                user.getName().isEmpty() ||
                user.getPassword() == null ||
                user.getPassword().isEmpty()) {
            return resultMap.failedMap(user, "name or password is empty", ErrCode.ARGS_FAILED);
        }

        String plainPassword = user.getPassword();
        user.setPassword(encryptUtils.EncryptByMd5(user.getPassword(), Conf.md5slat));
        int ret = 0;
        try {
            ret = userMapper.addUser(user);
        } catch (Exception e) {
            logger.warn("failed add user: ", e);
        }

        if (ret != 1) {
            user.setPassword(plainPassword);
            return resultMap.failedMap(user, "failed add user into db", ErrCode.DB_FAILED);
        }

        return resultMap.success();
    }

    @RequestMapping("/del")
    public Object delUser(HttpServletRequest request, @RequestBody User user) {
        logger.info("/v1.0/user/del");

        if (user.getName() == null || user.getName().isEmpty()) {
            return resultMap.failedMap(user, "user is empty", ErrCode.ARGS_FAILED);
        }

        int ret = 0;
        try {
            ret = userMapper.delUserByName(user.getName());
        } catch (Exception e) {
            logger.warn("failed del user: ", e);
        }

        if (ret != 1) {
            return resultMap.failedMap(user, "failed del user", ErrCode.DB_FAILED);
        }

        return resultMap.success();
    }

    @RequestMapping("/login")
    public Object userLogin(HttpServletRequest request, @RequestBody User user) throws NoSuchAlgorithmException {
        logger.info("/v1.0/user/login");

        if (user.getName() == null ||
                user.getName().isEmpty() ||
                user.getPassword() == null ||
                user.getPassword().isEmpty()) {
            return resultMap.failedMap(user, "name or password is empty", ErrCode.ARGS_FAILED);
        }

        String encryptPassword = encryptUtils.EncryptByMd5(user.getPassword(), Conf.md5slat);

        List<User> dbUsers = userMapper.getUserPasswordByName(user.getName());
        if (dbUsers == null || dbUsers.size() != 1) {
            return resultMap.failedMap(user, "name or password is false", ErrCode.DB_FAILED);
        }

        if (!encryptPassword.equals(dbUsers.get(0).getPassword())) {
            return resultMap.failedMap(user, "name or password is false", ErrCode.DB_FAILED);
        }

        return resultMap.success();
    }

    @RequestMapping("/update_password")
    public Object userUpdatePassword(HttpServletRequest request, @RequestBody UserPassword userPassword) throws NoSuchAlgorithmException {
        logger.info("/v1.0/user/update_password");

        if (userPassword.getName() == null || userPassword.getName().isEmpty() ||
                userPassword.getOldPassword() == null || userPassword.getOldPassword().isEmpty() ||
                userPassword.getNewPassword() == null || userPassword.getNewPassword().isEmpty()) {
            return resultMap.failedMap(userPassword, "name or password is empty", ErrCode.ARGS_FAILED);
        }

        String encryptOldPassword = encryptUtils.EncryptByMd5(userPassword.getOldPassword(), Conf.md5slat);
        String encryptNewPassword = encryptUtils.EncryptByMd5(userPassword.getNewPassword(), Conf.md5slat);
        int ret = userMapper.updateUserPassword(userPassword.getName(), encryptOldPassword, encryptNewPassword);
        if (ret != 1) {
            return resultMap.failedMap(userPassword, "failed update user password", ErrCode.DB_FAILED);
        }

        return resultMap.success();
    }

}
