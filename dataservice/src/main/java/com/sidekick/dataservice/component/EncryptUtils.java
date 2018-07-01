package com.sidekick.dataservice.component;

import org.springframework.stereotype.Component;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

@Component
public class EncryptUtils {

    public String EncryptByMd5(String plainText, String salt) throws NoSuchAlgorithmException {
        return md5(md5(plainText), md5(salt));
    }

    private String md5(String plainText) throws NoSuchAlgorithmException {
        return md5(plainText, null);
    }

    private String md5(String plainText, String salt)
            throws NoSuchAlgorithmException {
        MessageDigest md = MessageDigest.getInstance("MD5");

        if (salt != null) {
            md.update(salt.getBytes());
        }
        md.update(plainText.getBytes());

        byte byteData[] = md.digest();

        StringBuffer sb = new StringBuffer();
        for (int i = 0; i < byteData.length; i++) {
            sb.append(Integer.toString((byteData[i] & 0xff) + 0x100, 16)
                    .substring(1));
        }
        return sb.toString();
    }
}
