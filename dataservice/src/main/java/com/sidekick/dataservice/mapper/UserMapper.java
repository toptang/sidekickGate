package com.sidekick.dataservice.mapper;

import com.sidekick.dataservice.model.User;
import org.apache.ibatis.annotations.*;

import java.util.List;

@Mapper
public interface UserMapper {

    // 查询所有用户，所有信息
    @Select({"select id, name, password " +
            " from sidekick_user"})
    @Results(id = "userQueryMap", value = {
            @Result(property = "id", column = "id"),
            @Result(property = "name", column = "name"),
            @Result(property = "password", column = "password"),
    })
    List<User> getAllUsers();

    // 查询用户密码
    @Select({"select id, name, password " +
            " from sidekick_user " +
            " where name = #{name}"})
    @ResultMap("userQueryMap")
    List<User> getUserPasswordByName(@Param("name") String name);

    // 新建用户
    @Insert({"insert into sidekick_user (name, password) " +
            " values (#{user.name}, #{user.password})"})
    @Options(useGeneratedKeys = true, keyProperty = "user.id")
    int addUser(@Param("user") User user);

    // 删除用户
    @Delete({"delete from sidekick_user " +
            " where name=#{name}"})
    int delUserByName(@Param("name") String name);

    // 更新用户密码
    @Update({"update sidekick_user " +
            " set password=#{new_password} " +
            " where name=#{name} and password=#{old_password}"})
    int updateUserPassword(@Param("name") String name, @Param("old_password") String oldPassword, @Param("new_password") String newPassword);

}
