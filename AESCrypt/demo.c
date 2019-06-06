#include <stdio.h>
#include <pbt.h>
#include <ctype.h>

void aes(char*, char*, char*, int);
void aes_detail(int[4][4], int[4][4], int);
void subBytes(int [4][4], int);
void shiftRows(int [4][4], int);
void mixColumns(int [4][4], int);
void addRoundKey(int [4][4], int[4][4]);
int aes_multiple(int, int);
void keyExpansion(int key[4][4], int w[11][4][4]);
int c2i(char );

/**
 * S盒
 */
static const int S_BOX[16][16] = { 0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
                               0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
                               0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
                               0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
                               0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
                               0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
                               0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
                               0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
                               0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
                               0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
                               0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
                               0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
                               0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
                               0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
                               0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
                               0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16 };

/**
 * 逆S盒
 */
static const int INVERSE_S_BOX[16][16] = { 0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb,
                                0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb,
                                0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e,
                                0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25,
                                0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92,
                                0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84,
                                0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06,
                                0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b,
                                0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73,
                                0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e,
                                0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b,
                                0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4,
                                0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f,
                                0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef,
                                0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61,
                                0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d };
int RC[10] = {0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36};
int main(){
    int method = 0;//1表示加密， 0表示解密
    //待加密/解密文件存放路径
    char * source_path = "texten.txt";
    // 加密/解密后文件存放路径
    char *des_path = "text.txt";
    // 32位16进制密钥
    char * password = "0f1571c947d9e8590cb7add6af7f6798";
    aes(source_path, des_path, password, method);
    printf("success!!!!!!!!!!");

}

void aes(char* source_path, char* des_path, char* password, int method){
    //将密钥转换成4*4数组
    int p[4][4];
    for (int m = 0; m < 4; ++m) {
        for (int i = 0; i < 4; ++i) {
            int indx = 4 * i + m;
            p[i][m] = 16 * c2i(password[indx]) + c2i(password[indx + 1]);
        }
    }

    FILE *file = fopen(source_path, "r"); //获取文件的指针
    fseek(file, 0, SEEK_END); //移动文件的指针到文件结尾
    int len = ftell(file); //获取文件的长度
    rewind(file); //将文件指针移动回文件开始
    // 如果文件长度不是128位（16字节）的整数倍，则补齐
    int size = len;
    if (len % 16 != 0) {
        size = (len / 16 + 1) * 16;
    }
    unsigned char content[size];
    //读取文件内容赋值给content
    fread(content, 1, len, file);
    for (int j = len; j < size; ++j) {
        content[j] = 0;
    }
    fclose(file);
    //存储结果
    unsigned char encry[size];
    //将文件转换成16字节的int型数组加密、解密
    for (int i = 0; i < size / 16; ++i) {

        int content_to_int[4][4];
        for (int j = 0; j < 4; ++j) {
            for (int k = 0; k < 4; ++k) {
                content_to_int[j][k] = content[j * 4 + k + 16 * i];
            }
        }
        aes_detail(content_to_int, p, method);
        for (int j = 0; j < 4; ++j) {
            for (int k = 0; k < 4; ++k) {
                encry[j * 4 + k + 16 * i] = content_to_int[j][k];
            }
        }
    }
    FILE *file1 = fopen(des_path, "w");
    fwrite(encry, size, 1, file1);
    fflush(file1);
    fclose(file1);

}


void aes_detail(int content[4][4],  int password[4][4], int encode){
    int p[11][4][4];
    keyExpansion(password, p);

    if (encode) {
        addRoundKey(content, p[0]);
        for (int i = 1; i <= 10; ++i) {
            subBytes(content, encode);
            shiftRows(content, encode);
            if (i != 10) {
                mixColumns(content, encode);
            }

            addRoundKey(content, p[i]);
        }
    }else {
        addRoundKey(content, p[10]);
        for (int i = 9; i >= 0; --i) {
            shiftRows(content, encode);
            subBytes(content, encode);
            addRoundKey(content, p[i]);
            if (i != 0) {
                mixColumns(content, encode);
            }
        }
    }
}

void subBytes(int a[4][4], int encode){
    // encode 为1 代表字节替代，为0代表逆向字节替代
    for (int i = 0; i < 4; ++i) {
        for (int j = 0; j < 4; ++j) {
            int temp = a[i][j];
            int row = temp / 16;
            int column = temp % 16;
            if (encode)
                a[i][j] = S_BOX[row][column];
            else
                a[i][j] = INVERSE_S_BOX[row][column];
        }
    }
}

void shiftRows(int a[4][4], int encode){
    //encode 为1代表行移位，为0代表逆向行移位
    for (int i = 0; i < 4; ++i) {
        for (int j = 0; j < i; ++j) {
            if (encode) {
                int temp = a[i][0];
                a[i][0] = a[i][1];
                a[i][1] = a[i][2];
                a[i][2] = a[i][3];
                a[i][3] = temp;
            } else{
                int temp = a[i][3];
                a[i][3] = a[i][2];
                a[i][2] = a[i][1];
                a[i][1] = a[i][0];
                a[i][0] = temp;
            }
        }
    }
}

void mixColumns(int a[4][4], int encode){
    //encode 为1代表列混淆，为0代表逆向列混淆
    for (int i = 0; i < 4; ++i) {
        int temp0 = a[0][i];
        int temp1 = a[1][i];
        int temp2 = a[2][i];
        int temp3 = a[3][i];
        if (encode) {
            a[0][i] = aes_multiple(temp0, 2) ^ aes_multiple(temp1, 3) ^ temp2 ^ temp3;
            a[1][i] = temp0 ^ (aes_multiple(temp1, 2)) ^ (temp2 ^ aes_multiple(temp2, 2)) ^ temp3;
            a[2][i] = temp0 ^ temp1 ^ (aes_multiple(temp2, 2)) ^ (temp3 ^ aes_multiple(temp3, 2));
            a[3][i] = temp0 ^ (aes_multiple(temp0, 2)) ^ temp1 ^ temp2 ^ aes_multiple(temp3, 2);
        }else{
            a[0][i] = aes_multiple(temp0, 14) ^ aes_multiple(temp1, 11) ^ aes_multiple(temp2, 13) ^ aes_multiple(temp3, 9);
            a[1][i] = aes_multiple(temp0, 9) ^ aes_multiple(temp1, 14) ^ aes_multiple(temp2, 11) ^ aes_multiple(temp3, 13);
            a[2][i] = aes_multiple(temp0, 13) ^ aes_multiple(temp1, 9) ^ aes_multiple(temp2, 14) ^ aes_multiple(temp3, 11);
            a[3][i] = aes_multiple(temp0, 11) ^ aes_multiple(temp1, 13) ^ aes_multiple(temp2, 9) ^ aes_multiple(temp3, 14);
        }
    }
}

void addRoundKey(int a[4][4], int k[4][4]){
    // 由于用w[11][4][4]表示W[44]导致行列转置，所以在进行异或操作的时候应该是a[i，j] 异或 k[j,i]
    for (int i = 0; i < 4; ++i) {
        for (int j = 0; j < 4; ++j) {
            a[i][j] = a[i][j] ^ k[j][i];
        }
    }
}

//AES乘法计算
int aes_multiple(int a, int le){
    int thr = le & 0x8;
    int sec = le & 0x4;
    int fir = le & 0x2;
    int fir_mod = le % 2;
    int result = 0;
    if (thr){
        int b = a;
        for (int i = 1; i <=3 ; ++i) {
            b = b<<1;
            if (b >= 256)
                b = b ^ 0x11b;
        }
        b = b % 256;
        result = result ^ b;
    }
    if (sec){
        int b = a;
        for (int i = 1; i <=2 ; ++i) {
            b = b<<1;
            if (b >= 256)
                b = b ^ 0x11b;
        }
        b = b % 256;
        result = result ^ b;
    }
    if (fir){
        int b = a << 1;
        if (b >= 256)
            b = b ^ 0x11b;
        b = b % 256;
        result = result ^ b;
    }
    if (fir_mod)
        result = result ^ a;
    return  result;
}

void keyExpansion(int key[4][4], int w[11][4][4]){
    for (int i = 0; i < 4; ++i) {
        for (int j = 0; j < 4; ++j) {
            w[0][i][j] = key[j][i];
        }
    }
    for (int i = 1; i < 11; ++i){
        for (int j = 0; j < 4; ++j) {
            int temp[4];
            if (j == 0){
                temp[0] = w[i-1][3][1];
                temp[1] = w[i-1][3][2];
                temp[2] = w[i-1][3][3];
                temp[3] = w[i-1][3][0];
                for (int k = 0; k < 4; ++k) {
                    int m = temp[k];
                    int row = m / 16;
                    int column = m % 16;
                    temp[k] = S_BOX[row][column];
                    if (k == 0){
                        temp[k] = temp[k] ^ RC[i-1];
                    }
                }
            } else{
                temp[0] = w[i][j-1][0];
                temp[1] = w[i][j-1][1];
                temp[2] = w[i][j-1][2];
                temp[3] = w[i][j-1][3];
            }
            for (int l = 0; l < 4; ++l) {

                w[i][j][l] = w[i-1][j][l] ^ temp[l];
            }

        }
    }

}

//将字符转换为数值
int c2i(char ch) {
    // 如果是数字，则用数字的ASCII码减去48, 如果ch = '2' ,则 '2' - 48 = 2
    if(isdigit(ch))
        return ch - 48;

    // 如果是字母，但不是A~F,a~f则返回
    if( ch < 'A' || (ch > 'F' && ch < 'a') || ch > 'z' )
        return -1;

    // 如果是大写字母，则用数字的ASCII码减去55, 如果ch = 'A' ,则 'A' - 55 = 10
    // 如果是小写字母，则用数字的ASCII码减去87, 如果ch = 'a' ,则 'a' - 87 = 10
    if(isalpha(ch))
        return isupper(ch) ? ch - 55 : ch - 87;

    return -1;
}
