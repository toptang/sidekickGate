package com.sidekick.dataservice.controller;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequestMapping("/api/v1.0/user")
@RestController
public class HelloController {
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

//    @Autowired
//    private UserMapper userMapper;
//
//    @Autowired
//    private ResultMap resultMap;
//
//    @RequestMapping("/all")
//    public Object getAllUsers() {
//        logger.info("/api/v1.0/users/all");
//        List<User> users = userMapper.getAllUsers();
//        if (users != null) {
//            return resultMap.successMap(users);
//        }
//        return resultMap.failedMap(null, "failed query users", 1);
//    }
//
//    @RequestMapping("/code")
//    public Object getUserByCode(HttpServletRequest request, @RequestBody Map<String,Object> map) {
//        logger.info("/api/v1.0/users/code");
//
//        String code = (String)map.get("code");
//        List<User> users = userMapper.getUserByCode(code);
//        if (users != null && users.size() == 1) {
//            return resultMap.successMap(users.get(0));
//        }
//        return resultMap.failedMap(map, "failed query user by code", 1);
//    }
}
