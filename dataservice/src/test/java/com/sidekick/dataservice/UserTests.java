package com.sidekick.dataservice;

import com.sidekick.dataservice.component.EncryptUtils;
import com.sidekick.dataservice.enumeration.Conf;
import com.sidekick.dataservice.mapper.UserMapper;
import com.sidekick.dataservice.model.User;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit.jupiter.SpringExtension;

import java.io.UnsupportedEncodingException;
import java.security.NoSuchAlgorithmException;
import java.util.List;

@ExtendWith(SpringExtension.class)
@SpringBootTest
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class UserTests {
    @Autowired
    private UserMapper userMapper;

    @Autowired
    private EncryptUtils encryptUtils;

    private String[] usernames = new String[]{
            "foo", "bar"
    };
    private String passwd = "foobar!@#456&*";

    @BeforeAll
    public void addUsers() throws UnsupportedEncodingException, NoSuchAlgorithmException {
        String encryptPasswd = encryptUtils.EncryptByMd5(passwd, Conf.md5slat);

        for (String username : usernames) {
            User user = new User();
            user.setName(username);
            user.setPassword(encryptPasswd);
            int ret = userMapper.addUser(user);
            Assertions.assertEquals(ret, 1);
        }
    }

    @AfterAll
    public void delUsers() {
        for (String username : usernames) {
            int ret = userMapper.delUserByName(username);
            Assertions.assertEquals(ret, 1);
        }
    }

    @Test
    public void testUserPasswd() throws NoSuchAlgorithmException {
        String encryptPasswd = encryptUtils.EncryptByMd5(passwd, Conf.md5slat);

        String newPasswd = "xxx123";
        String encryptNewPasswd = encryptUtils.EncryptByMd5(newPasswd, Conf.md5slat);

        for (String username : usernames) {
            List<User> users = userMapper.getUserPasswordByName(username);
            Assertions.assertNotNull(users);
            Assertions.assertEquals(users.size(), 1);
            User user = users.get(0);
            Assertions.assertEquals(user.getName(), username);
            Assertions.assertEquals(user.getPassword(), encryptPasswd);

            int ret = userMapper.updateUserPassword(username, encryptPasswd, encryptNewPasswd);
            Assertions.assertEquals(ret, 1);

            users = userMapper.getUserPasswordByName(username);
            Assertions.assertNotNull(users);
            Assertions.assertEquals(users.size(), 1);
            user = users.get(0);
            Assertions.assertEquals(user.getName(), username);
            Assertions.assertEquals(user.getPassword(), encryptNewPasswd);
        }
    }

}
