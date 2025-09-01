const crypto = require('crypto');

// MD5加密函数
function md5(str) {
    return crypto.createHash('md5').update(str).digest('hex');
}

// 根据原始代码逆向的签名生成函数
function S(e, t, n, r) {
    try {
        var o, l;

        // 配置对象 - 对应原代码中的p
        const p = {
            apiKey: 'e98ce2565b09ecc0',
            apiSecret: 'b28efc98ae90a878'
        };

        // 获取API密钥和密钥 - 对应原代码逻辑
        if (n) {
            o = p["".concat(n, "ApiKey")];
            l = p["".concat(n, "ApiSecret")];
        } else {
            o = p.apiKey;
            l = p.apiSecret;
        }

        var d, h = new Date().getTime();
        var centerId = '50030001';
        var tenantId = '82';

        // 构建参数对象 - 对应原代码的 (0, i.default)({}, {...}, t)
        var y = Object.assign({}, {
            apiKey: o,
            timestamp: h,
            channelId: '11'
        }, t);

        // 添加centerId - 对应原代码逻辑：!g || y.centerId || 0 === y.centerId || (r || {}).noCenterId || (y.centerId = g)
        if (!(!centerId || y.centerId || y.centerId === 0 || (r || {}).noCenterId)) {
            y.centerId = centerId;
        }

        // 添加tenantId - 对应原代码逻辑：v && (y.tenantId = v)
        if (tenantId) {
            y.tenantId = tenantId;
        }

        // 转换为键值对数组 - 对应原代码的 (0, a.default)(y, function(e, t) {...})
        var m = Object.keys(y).map(function (key) {
            return {
                key: key,
                value: y[key]
            };
        });

        // 按key排序 - 对应原代码的 (0, u.default)(m, "key")
        m = m.sort(function (a, b) {
            return a.key.localeCompare(b.key);
        });

        // 拼接参数字符串 - 对应原代码的 (0, c.default)(m, function(e, t) {...}, "")
        d = m.reduce(function (acc, item) {
            return "".concat(acc).concat(item.key, "=").concat(item.value);
        }, "");

        // 生成待签名字符串并编码 - 对应原代码：var _ = encodeURIComponent(e + d + l);
        var _ = encodeURIComponent(e + d + l);

        // 替换特殊字符 - 严格按照原代码逻辑
        if (_.indexOf("(") >= 0) {
            _ = _.replace(/\(/g, "%28");
        }
        if (_.indexOf(")") >= 0) {
            _ = _.replace(/\)/g, "%29");
        }
        if (_.indexOf("'") >= 0) {
            _ = _.replace(/'/g, "%27");
        }
        if (_.indexOf("!") >= 0) {
            _ = _.replace(/!/g, "%21");
        }
        if (_.indexOf("~") >= 0) {
            _ = _.replace(/~/g, "%7E");
        }

        // MD5加密 - 对应原代码的 d = (0, s.default)(_)
        d = md5(_);

        // 返回包含签名的参数对象 - 对应原代码的 (0, i.default)({}, y, {sign: d})
        return Object.assign({}, y, { sign: d });

    } catch (b) {
        console.log(b);
    }
}

// 通用签名生成函数
function generateSignature(apiPath, params = {}, options = {}) {
    return S(apiPath, params, options.prefix, options);
}

// 将generateSignature生成的对象转换为URL参数字符串
function toUrlParams(signatureObject) {
    if (!signatureObject || typeof signatureObject !== 'object') {
        return '';
    }

    const params = [];

    // 遍历对象的所有属性
    for (const key in signatureObject) {
        if (signatureObject.hasOwnProperty(key)) {
            const value = signatureObject[key];
            // URL编码键和值
            const encodedKey = encodeURIComponent(key);
            const encodedValue = encodeURIComponent(value);
            params.push(`${encodedKey}=${encodedValue}`);
        }
    }

    // 用&连接所有参数
    return params.join('&');
}

// 生成完整的URL（包含路径和参数）
function generateFullUrl(baseUrl, apiPath, params = {}, options = {}) {
    // 生成签名对象
    const signatureObject = generateSignature(apiPath, params, options);

    // 转换为URL参数
    const urlParams = toUrlParams(signatureObject);

    // 构建完整URL
    const fullUrl = `${baseUrl}${apiPath}?${urlParams}`;

    return {
        url: fullUrl,
        params: signatureObject,
        queryString: urlParams
    };
}

// 命令行参数处理
function parseCommandLineArgs() {
    const args = process.argv.slice(2); // 去掉 node 和脚本名
    const options = {
        method: 'getProducts', // 默认方法
        params: {}
    };

    for (let i = 0; i < args.length; i++) {
        const arg = args[i];

        if (arg === '--method' || arg === '-m') {
            // 指定签名方法
            options.method = args[i + 1];
            i++; // 跳过下一个参数
        } else if (arg === '--serviceId' || arg === '-s') {
            // 指定服务ID
            options.params.serviceId = args[i + 1];
            i++;
        } else if (arg === '--venueId' || arg === '-v') {
            // 指定场馆ID
            options.params.venueId = args[i + 1];
            i++;
        } else if (arg === '--timestamp' || arg === '-t') {
            // 指定时间戳
            options.timestamp = parseInt(args[i + 1]);
            i++;
        } else if (arg === '--help' || arg === '-h') {
            // 显示帮助信息
            showHelp();
            process.exit(0);
        } else if (arg.startsWith('--')) {
            // 处理其他自定义参数 --key=value 或 --key value 格式
            if (arg.includes('=')) {
                // --key=value 格式
                const [key, value] = arg.substring(2).split('=');
                if (value !== undefined) {
                    options.params[key] = value;
                }
            } else {
                // --key value 格式
                const key = arg.substring(2);
                if (i + 1 < args.length && !args[i + 1].startsWith('-')) {
                    options.params[key] = args[i + 1];
                    i++; // 跳过下一个参数
                } else {
                    // 如果没有值，设置为空字符串
                    options.params[key] = '';
                }
            }
        }
    }

    return options;
}

// 显示帮助信息
function showHelp() {
    console.log(`
签名生成器使用说明:

基本用法:
  node signature_generator.js [选项]

选项:
  -m, --method <method>     指定签名方法 (默认: getProducts)
                           可选值: getProducts, getProductsGeneration, fieldList, newOrder, custom
  
  -s, --serviceId <id>      指定服务ID (默认: 1001)
  -v, --venueId <id>        指定场馆ID (默认: 5003000101)
  -t, --timestamp <ts>      指定时间戳 (默认: 当前时间)
  
  --<key>=<value>          添加自定义参数
  
  -h, --help               显示此帮助信息

示例:
  # 使用默认参数生成getProducts签名
  node signature_generator.js
  
  # 指定方法和参数
  node signature_generator.js --method getProductsGeneration --serviceId 1002 --venueId 5003000102
  
  # 生成fieldList签名
  node signature_generator.js --method fieldList --netUserId 2025082802482655 --day 20250830
  
  # 生成newOrder签名
  node signature_generator.js --method newOrder --fieldInfo "f3a74ed8bf6efbe143e33c0ba1cf9e26,fade00ab4725896da1840e9c0125dc0f"
  
  # 使用自定义参数
  node signature_generator.js --method custom --customParam=value --anotherParam=123
  
  # 指定时间戳
  node signature_generator.js --timestamp 1756366123577
`);
}

// 不同的签名方法
const signatureMethods = {
    // 原始getProducts方法
    getProducts: function (options) {
        const apiPath = '/aisports-api/api/product/getProducts';
        const params = {
            serviceId: options.params.serviceId || '1001',
            venueId: options.params.venueId || '5003000101'
        };

        // 生成签名对象
        const signatureObject = generateSignature(apiPath, params, options);

        // 转换为URL参数
        return [toUrlParams(signatureObject)];
    },

    // fieldList方法 - 根据参考URL生成场地列表签名
    fieldList: function (options) {
        if (!options.params.day) {
            return
        }
        const apiPath = '/aisports-api/wechatAPI/venue/fieldList';
        const params = {
            netUserId: options.params.netUserId || '2025082802482655',
            venueId: options.params.venueId || '5003000101',
            serviceId: options.params.serviceId || '1002',
            day: options.params.day,
            selectByfullTag: options.params.selectByfullTag || '0',
            fieldType: options.params.fieldType || '1841',
        };

        // 生成签名对象
        const signatureObject = generateSignature(apiPath, params, options);

        // 转换为URL参数
        return [toUrlParams(signatureObject)];
    },

    // newOrder方法 - 根据参考URL生成新订单签名
    newOrder: function (options) {
        if (!(options.params.fieldInfo && options.params.day)) {
            return
        }
        const apiPath = '/aisports-api/wechatAPI/order/newOrder';
        const params = {
            serviceId: options.params.serviceId || '1002',
            day: options.params.day,
            fieldType: options.params.fieldType || '1841',
            fieldInfo: options.params.fieldInfo,
            ticket: options.params.ticket || '',
            randStr: options.params.randStr || '',
            venueId: '5003000101',
            netUserId: options.params.netUserId || '2025082802482655'
        };

        // 生成签名对象
        const signatureObject = generateSignature(apiPath, params, options);

        // 转换为URL参数
        return [toUrlParams(signatureObject)];
    },

    // 自定义方法
    custom: function (options) {
        const apiPath = options.params.apiPath || '/aisports-api/api/custom/endpoint';
        const params = { ...options.params };
        delete params.apiPath; // 移除apiPath，因为它不是请求参数

        const baseUrl = 'https://web.xports.cn';
        const result = generateFullUrl(baseUrl, apiPath, params, options);

        console.log('方法: custom');
        console.log('API路径:', apiPath);
        console.log('参数:', params);
        console.log('查询字符串:', result.queryString);
        console.log('完整URL:', result.url);

        return [result];
    }
};

// 主执行函数
function main() {
    try {
        const options = parseCommandLineArgs();
        // 根据方法名选择对应的签名方法
        const method = signatureMethods[options.method];
        if (!method) {
            console.error(`错误: 未知的签名方法 "${options.method}"`);
            console.error('可用方法:', Object.keys(signatureMethods).join(', '));
            process.exit(1);
        }

        // 执行选定的方法
        method(options).forEach(element => {
            console.log(element);
        });
    } catch (error) {
        console.error('执行错误:', error.message);
        process.exit(1);
    }
}

// 保留原有的getProductsGeneration函数用于向后兼容
function getProductsGeneration() {
    const apiPath = '/aisports-api/api/product/getProducts';
    const params = {
        serviceId: '1001',
        venueId: '5003000101'
    };

    const baseUrl = 'https://web.xports.cn';
    const fullUrlResult = generateFullUrl(baseUrl, apiPath, params);

    return fullUrlResult;
}

// 导出函数
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        S,
        generateSignature,
        toUrlParams,
        generateFullUrl,
        signatureMethods,
        parseCommandLineArgs,
        main,
        getProductsGeneration
    };
}

// 如果是直接运行此脚本，则执行main函数
if (require.main === module) {
    main();
}

