"""
Minimal Payment API: V1 → V2 Upgrade Demo
Simple example showing API versioning strategy
"""
from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

# Simple in-memory storage
payments = {}
payment_id = 1000


# v1: Simple
@app.route('/api/v1/payments', methods=['POST'])
def v1_create():
    """V1: Create payment (simple)"""
    global payment_id
    data = request.get_json()
    
    if not data or 'amount' not in data:
        return {'error': 'Missing amount'}, 400
    
    payment_id += 1
    payment = {
        'id': payment_id,
        'amount': data['amount'],
        'currency': data.get('currency', 'USD'),
        'status': 'pending'
    }
    payments[payment_id] = payment
    
    return {
        'payment_id': payment_id,
        'amount': payment['amount'],
        'status': 'pending'
    }, 201


@app.route('/api/v1/payments/<int:pid>', methods=['GET'])
def v1_get(pid):
    """V1: Get payment"""
    if pid not in payments:
        return {'error': 'Not found'}, 404
    
    p = payments[pid]
    return {
        'payment_id': p['id'],
        'amount': p['amount'],
        'status': p['status']
    }, 200


@app.route('/api/v1/payments/<int:pid>', methods=['PUT'])
def v1_update(pid):
    """V1: Update payment status"""
    if pid not in payments:
        return {'error': 'Not found'}, 404
    
    data = request.get_json()
    payments[pid]['status'] = data.get('status', 'pending')
    
    return {
        'payment_id': pid,
        'status': payments[pid]['status']
    }, 200


# v2: Enhanced

@app.route('/api/v2/payments', methods=['POST'])
def v2_create():
    """V2: Create payment (enhanced with customer_id, metadata)"""
    global payment_id
    data = request.get_json()
    
    if not data or 'amount' not in data or 'customer_id' not in data:
        return {'error': 'Missing amount or customer_id'}, 400
    
    payment_id += 1
    payment = {
        'id': payment_id,
        'amount': data['amount'],
        'currency': data.get('currency', 'USD'),
        'customer_id': data['customer_id'],
        'description': data.get('description', ''),
        'metadata': data.get('metadata', {}),
        'status': 'pending',
        'created_at': datetime.now().isoformat(),
        'updated_at': datetime.now().isoformat()
    }
    payments[payment_id] = payment
    
    return {
        'status': 'success',
        'payment': {
            'id': payment_id,
            'amount': payment['amount'],
            'customer_id': payment['customer_id'],
            'status': 'pending',
            'created_at': payment['created_at']
        }
    }, 201


@app.route('/api/v2/payments/<int:pid>', methods=['GET'])
def v2_get(pid):
    """V2: Get payment with all details"""
    if pid not in payments:
        return {'error': 'Not found'}, 404
    
    p = payments[pid]
    return {
        'status': 'success',
        'payment': {
            'id': p['id'],
            'amount': p['amount'],
            'currency': p['currency'],
            'customer_id': p.get('customer_id'),
            'description': p.get('description', ''),
            'metadata': p.get('metadata', {}),
            'status': p['status'],
            'created_at': p.get('created_at'),
            'updated_at': p.get('updated_at')
        }
    }, 200


@app.route('/api/v2/payments/<int:pid>', methods=['PUT'])
def v2_update(pid):
    """V2: Update payment (status, metadata, description)"""
    if pid not in payments:
        return {'error': 'Not found'}, 404
    
    data = request.get_json()
    
    if 'status' in data:
        payments[pid]['status'] = data['status']
    if 'metadata' in data:
        payments[pid]['metadata'] = {**payments[pid].get('metadata', {}), **data['metadata']}
    if 'description' in data:
        payments[pid]['description'] = data['description']
    
    payments[pid]['updated_at'] = datetime.now().isoformat()
    
    return {
        'status': 'success',
        'payment': payments[pid]
    }, 200


@app.route('/api/v2/payments', methods=['GET'])
def v2_list():
    """V2: List payments (NEW feature)"""
    customer_id = request.args.get('customer_id')
    
    result = []
    for p in payments.values():
        if customer_id and p.get('customer_id') != customer_id:
            continue
        result.append(p)
    
    return {
        'status': 'success',
        'count': len(result),
        'payments': result
    }, 200



@app.route('/api/health', methods=['GET'])
def health():
    """Health check"""
    return {
        'status': 'healthy',
        'versions': ['v1', 'v2']
    }, 200


@app.route('/', methods=['GET'])
def index():
    """API Info"""
    return {
        'name': 'Payment API',
        'endpoints': {
            'v1': [
                'POST /api/v1/payments',
                'GET /api/v1/payments/<id>',
                'PUT /api/v1/payments/<id>'
            ],
            'v2': [
                'POST /api/v2/payments',
                'GET /api/v2/payments/<id>',
                'GET /api/v2/payments?customer_id=X',
                'PUT /api/v2/payments/<id>'
            ]
        }
    }, 200


if __name__ == '__main__':
    app.run(debug=True, port=5000)
