const hooks = require('hooks');
const axios = require('axios');

// Base URL configuration
const BASE_URL = 'http://127.0.0.1:8080';

// Global değişkenler
let authToken = '';
let testUserId = null;
let createdTaskId = null;
let testUserEmail = '';
let testUserPassword = 'test123456';

console.log('🚀 Starting Dredd API Tests...');

// Test ortamını hazırla
hooks.beforeAll(async (transactions, done) => {
  console.log('📋 Setting up test environment...');
  
  const timestamp = Date.now();
  const testEmail = `dredd_test_${timestamp}@test.com`;
  const testUsername = `dredd_test_${timestamp}`;
  const testPassword = 'test123456';

  try {
    // Test kullanıcısını kaydet
    const registerResponse = await axios.post(`${BASE_URL}/register`, {
      username: testUsername,
      email: testEmail, 
      password: testPassword
    });
    
    // Global olarak sakla
    testUserEmail = testEmail;
    
    if (registerResponse.status === 201) {
      console.log('✅ User registered successfully');
    }
  } catch (error) {
    if (error.response?.status === 400) {
      console.log('ℹ️ User already exists, continuing...');
    } else {
      console.log('❌ Registration failed:', error.message);
    }
  }

  try {
    // Kullanıcıyı login et
    const loginResponse = await axios.post(`${BASE_URL}/login`, {
      email: testEmail,
      password: testPassword
    });
    
    authToken = loginResponse.data.token;
    testUserId = loginResponse.data.user.id;
    console.log('✅ Authentication successful, token obtained');

    // Test task'ı oluştur
    const createTaskResponse = await axios.post(`${BASE_URL}/tasks`, {
      title: 'Test Task for Dredd',
      description: 'This is a test task for Dredd testing',
      status: 'pending',
      priority: 'high'
    }, {
      headers: { Authorization: `Bearer ${authToken}` }
    });

    if (createTaskResponse.status === 201) {
      createdTaskId = createTaskResponse.data.id;
      console.log('✅ Test task created with ID:', createdTaskId);
    }

  } catch (error) {
    console.log('❌ Authentication failed:', error.message);
  }

  console.log('✅ Test environment ready');
  done();
});

// Her test öncesi
hooks.beforeEach((transaction, done) => {
  const testName = `${transaction.request.method} ${transaction.fullPath}`;
  const expectedCode = transaction.expected.statusCode;
  console.log(`📊 ${testName} - ${expectedCode}`);
  console.log(`🔍 Hook name: ${transaction.name}`);
  done();
});

// Authorization header'ı gerekli testlere ekle
hooks.beforeEach((transaction, done) => {
  // Auth gerektiren endpoint'ler
  const protectedPaths = ['/tasks', '/logout'];
  const isProtected = protectedPaths.some(path => transaction.fullPath.includes(path)) && 
                     !transaction.fullPath.includes('/public');

  if (isProtected && authToken) {
    transaction.request.headers['Authorization'] = `Bearer ${authToken}`;
  }
  done();
});

// Register testleri için özel işlem
hooks.before('/register > Register a new user > 201 > application/json', (transaction, done) => {
  // Yeni bir unique kullanıcı adı oluştur
  const timestamp = Date.now();
  transaction.request.body = JSON.stringify({
    username: `test_user_${timestamp}`,
    email: `test_${timestamp}@example.com`,
    password: 'password123'
  });
  done();
});

// Login 401 testi için yanlış şifre kullan
hooks.before('/login > Login user > 401 > application/json', (transaction, done) => {
  transaction.request.body = JSON.stringify({
    email: 'nonexistent@test.com',
    password: 'wrong_password'
  });
  done();
});

// Logout 401 testi için geçersiz token kullan
hooks.before('/logout > Logout user > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  done();
});

// Register 400 testi için invalid data kullan
hooks.before('/register > Register a new user > 400 > application/json', (transaction, done) => {
  // Eksik veya geçersiz veri gönder
  transaction.request.body = JSON.stringify({
    username: '', // Boş username
    email: 'invalid-email', // Geçersiz email format
    password: '123' // Çok kısa password
  });
  done();
});

// Login 200 testi için setup sırasında oluşturulan kullanıcıyı kullan
hooks.before('/login > Login user > 200 > application/json', (transaction, done) => {
  console.log('🔑 Login 200 test - global test user kullanılıyor');
  
  // Global test user credential'larını kullan  
  transaction.request.body = JSON.stringify({
    email: testUserEmail,
    password: testUserPassword
  });
  
  done();
});

// Tasks 401 testleri için geçersiz token kullan
hooks.before('/tasks > Get user tasks > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  done();
});

hooks.before('/tasks > Create new task > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  done();
});

// Tasks 400 testi için eksik veri gönder
hooks.before('/tasks > Create new task > 400 > application/json', (transaction, done) => {
  transaction.request.body = JSON.stringify({
    title: '' // Boş title göndererek 400 hatası oluştur
  });
  done();
});

// Task ID testleri için oluşturulan task ID'sini kullan
hooks.before('/tasks/{id} > Get task by ID > 200 > application/json', (transaction, done) => {
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

hooks.before('/tasks/{id} > Update task > 200 > application/json', (transaction, done) => {
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

hooks.before('/tasks/{id} > Delete task > 200 > application/json', (transaction, done) => {
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

// Task 401 testleri için geçersiz token kullan
hooks.before('/tasks/{id} > Get task by ID > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

hooks.before('/tasks/{id} > Update task > 400 > application/json', (transaction, done) => {
  transaction.request.body = JSON.stringify({
    title: '' // Boş title göndererek 400 hatası oluştur
  });
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

hooks.before('/tasks/{id} > Update task > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

hooks.before('/tasks/{id} > Delete task > 401 > application/json', (transaction, done) => {
  transaction.request.headers['Authorization'] = 'Bearer invalid_token_here';
  if (createdTaskId) {
    transaction.fullPath = transaction.fullPath.replace('/tasks/1', `/tasks/${createdTaskId}`);
    transaction.request.uri = transaction.request.uri.replace('/tasks/1', `/tasks/${createdTaskId}`);
  }
  done();
});

// 404 testleri için non-existent task ID kullan (değiştirme!)
hooks.before('/tasks/{id} > Get task by ID > 404 > application/json', (transaction, done) => {
  // Var olmayan task ID kullan (değiştirme yapma)
  done();
});

hooks.before('/tasks/{id} > Update task > 404 > application/json', (transaction, done) => {
  // Var olmayan task ID kullan (değiştirme yapma)
  done();
});

hooks.before('/tasks/{id} > Delete task > 404 > application/json', (transaction, done) => {
  // Var olmayan task ID kullan (değiştirme yapma)
  done();
});

// Test sonrası temizlik
hooks.afterAll(async (transactions, done) => {
  console.log('\n🏁 Dredd API Tests Completed');
  console.log('🧹 Cleaning up test environment...');
  
  try {
    // Test verilerini temizle
    if (createdTaskId && authToken) {
      await axios.delete(`${BASE_URL}/tasks/${createdTaskId}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      });
    }
  } catch (error) {
    // Sessizce devam et
  }

  console.log('✅ Cleanup completed');
  done();
});
